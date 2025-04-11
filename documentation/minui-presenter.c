#include <fcntl.h>
#include <getopt.h>
#include <msettings.h>
#include <parson/parson.h>
#include <pthread.h>
#include <signal.h>
#include <stdatomic.h>
#include <stdbool.h>
#include <stdio.h>
#include <sys/time.h>
#include <unistd.h>
#ifdef USE_SDL2
#include <SDL2/SDL_ttf.h>
#include <SDL2/SDL_image.h>
#else
#include <SDL/SDL_ttf.h>
#include <SDL/SDL_image.h>
#endif

#include "defines.h"
#include "api.h"
#include "utils.h"

SDL_Surface *screen = NULL;

#ifdef USE_SDL2
bool use_sdl2 = true;
#else
bool use_sdl2 = false;
#endif

pthread_mutex_t increment_item_list_index_lock = PTHREAD_MUTEX_INITIALIZER;
volatile sig_atomic_t increment_item_list_index = 0;

enum list_result_t
{
    ExitCodeSuccess = 0,
    ExitCodeError = 1,
    ExitCodeCancelButton = 2,
    ExitCodeMenuButton = 3,
    ExitCodeActionButton = 4,
    ExitCodeInactionButton = 5,
    ExitCodeStartButton = 6,
    ExitCodeParseError = 10,
    ExitCodeSerializeError = 11,
    ExitCodeTimeout = 124,
    ExitCodeKeyboardInterrupt = 130,
    ExitCodeSigterm = 143,
};
typedef int ExitCode;

// log_error logs a message to stderr for debugging purposes
void log_error(const char *msg)
{
    // Set stderr to unbuffered mode
    setvbuf(stderr, NULL, _IONBF, 0);
    fprintf(stderr, "%s\n", msg);
}

// log_info logs a message to stdout for debugging purposes
void log_info(const char *msg)
{
    // Set stdout to unbuffered mode
    setvbuf(stdout, NULL, _IONBF, 0);
    printf("%s\n", msg);
}

// Fonts holds the fonts for the list
struct Fonts
{
    // the size of the font to use for the list
    int size;
    // the large font to use for the list
    TTF_Font *large;
    // the medium font to use for the list
    TTF_Font *medium;
    // the small font to use for the list
    TTF_Font *small;

    // the path to the font to use for the list
    char *font_path;
};

enum MessageAlignment
{
    MessageAlignmentTop,
    MessageAlignmentMiddle,
    MessageAlignmentBottom,
};

struct Item
{
    // the background color to use for the list
    char *background_color;
    // path to the background image to use for the list
    char *background_image;
    // whether the background image exists
    bool image_exists;
    // the text to display
    char *text;
    // whether to show a pill around the text or not
    bool show_pill;
    // the alignment of the text
    enum MessageAlignment alignment;
};

// ItemsState holds the state of the list
struct ItemsState
{
    // array of display items
    struct Item *items;
    // number of items in the list
    size_t item_count;
    // index of currently selected item
    int selected;
};

// AppState holds the current state of the application
struct AppState
{
    // whether the screen needs to be redrawn
    int redraw;
    // whether the app should exit
    int quitting;
    // the exit code to return
    int exit_code;
    // the button to display on the Action button
    char action_button[1024];
    // whether to show the Action button
    bool action_show;
    // the text to display on the Action button
    char action_text[1024];
    // the background image to display
    char background_image[1024];
    // the background color to display
    char background_color[1024];
    // the button to display on the Confirm button
    char confirm_button[1024];
    // whether to show the Confirm button
    bool confirm_show;
    // the text to display on the Confirm button
    char confirm_text[1024];
    // the button to display on the Cancel button
    char cancel_button[1024];
    // whether to show the Cancel button
    bool cancel_show;
    // the text to display on the Cancel button
    char cancel_text[1024];
    // the button to display on the Inaction button
    char inaction_button[1024];
    // whether to show the Inaction button
    bool inaction_show;
    // the text to display on the Inaction button
    char inaction_text[1024];
    // the path to the JSON file
    char file[1024];
    // quit after last item
    bool quit_after_last_item;
    // whether to show the hardware group
    bool show_hardware_group;
    // whether to show the pill
    bool show_pill;
    // whether to show the time left
    bool show_time_left;
    // the seconds to display the message for before timing out
    int timeout_seconds;
    // the key to the items array in the JSON file
    char item_key[1024];
    // the start time of the presentation
    struct timeval start_time;
    // the fonts to use for the list
    struct Fonts fonts;
    // the display states
    struct ItemsState *items_state;
};

struct Message
{
    char message[1024];
    int width;
};

void strtrim(char *s)
{
    char *p = s;
    int l = strlen(p);

    if (l == 0)
    {
        return;
    }

    while (isspace(p[l - 1]))
    {
        p[--l] = 0;
    }

    while (*p && isspace(*p))
    {
        ++p;
        --l;
    }

    memmove(s, p, l + 1);
}

char *read_stdin()
{
    // Read all of stdin into a string
    char *stdin_contents = NULL;
    size_t stdin_size = 0;
    size_t stdin_used = 0;
    char buffer[4096];
    size_t bytes_read;

    while ((bytes_read = fread(buffer, 1, sizeof(buffer), stdin)) > 0)
    {
        if (stdin_contents == NULL)
        {
            stdin_size = bytes_read * 2;
            stdin_contents = malloc(stdin_size);
        }
        else if (stdin_used + bytes_read > stdin_size)
        {
            stdin_size *= 2;
            stdin_contents = realloc(stdin_contents, stdin_size);
        }

        memcpy(stdin_contents + stdin_used, buffer, bytes_read);
        stdin_used += bytes_read;
    }

    // Null terminate the string
    if (stdin_contents)
    {
        if (stdin_used == stdin_size)
        {
            stdin_contents = realloc(stdin_contents, stdin_size + 1);
        }
        stdin_contents[stdin_used] = '\0';
    }

    return stdin_contents;
}

// hydrate_display_states hydrates the display states from a file or stdin
struct ItemsState *ItemsState_New(const char *filename, const char *item_key, const char *default_background_image, const char *default_background_color, bool default_show_pill, enum MessageAlignment default_alignment)
{
    struct ItemsState *state = malloc(sizeof(struct ItemsState));

    JSON_Value *root_value;
    if (strcmp(filename, "-") == 0)
    {
        char *contents = read_stdin();
        if (contents == NULL)
        {
            log_error("Failed to read stdin");
            return NULL;
        }

        root_value = json_parse_string_with_comments(contents);
        free(contents);
    }
    else
    {
        root_value = json_parse_file_with_comments(filename);
    }

    if (root_value == NULL)
    {
        log_error("Failed to parse JSON file");
        return NULL;
    }

    JSON_Object *root_object = json_value_get_object(root_value);
    if (root_object == NULL)
    {
        json_value_free(root_value);
        return NULL;
    }

    JSON_Array *items = json_object_get_array(root_object, item_key);
    if (items == NULL)
    {
        json_value_free(root_value);
        return NULL;
    }

    size_t item_count = json_array_get_count(items);
    if (item_count == 0)
    {
        json_value_free(root_value);
        return NULL;
    }

    state->items = malloc(sizeof(struct Item) * item_count);

    for (size_t i = 0; i < item_count; i++)
    {
        JSON_Object *item = json_array_get_object(items, i);
        if (item == NULL)
        {
            char buff[1024];
            snprintf(buff, sizeof(buff), "Failed to get item %zu", i);
            log_error(buff);
            json_value_free(root_value);
            return NULL;
        }

        const char *text = json_object_get_string(item, "text");
        if (text == NULL)
        {
            char buff[1024];
            snprintf(buff, sizeof(buff), "Failed to get text for item %zu", i);
            log_error(buff);
            json_value_free(root_value);
            return NULL;
        }

        state->items[i].text = strdup(text);

        const char *background_image = json_object_get_string(item, "background_image");
        state->items[i].background_image = strdup(default_background_image);
        state->items[i].image_exists = default_background_image != NULL && access(default_background_image, F_OK) != -1;
        if (background_image != NULL)
        {
            state->items[i].background_image = strdup(background_image);
            state->items[i].image_exists = access(background_image, F_OK) != -1;
        }

        const char *background_color = json_object_get_string(item, "background_color");
        state->items[i].background_color = strdup(default_background_color);
        if (background_color != NULL)
        {
            state->items[i].background_color = strdup(background_color);
        }

        state->items[i].show_pill = default_show_pill;
        if (json_object_has_value(item, "show_pill"))
        {
            if (json_object_get_boolean(item, "show_pill") == 1)
            {
                state->items[i].show_pill = true;
            }
            else if (json_object_get_boolean(item, "show_pill") == 0)
            {
                state->items[i].show_pill = false;
            }
            else
            {
                char buff[1024];
                snprintf(buff, sizeof(buff), "Invalid show_pill value provided for item %zu", i);
                log_error(buff);
                json_value_free(root_value);
                return NULL;
            }
        }

        const char *alignment = json_object_get_string(item, "alignment");
        if (alignment == NULL)
        {
            state->items[i].alignment = default_alignment;
        }
        else if (strcmp(alignment, "top") == 0)
        {
            state->items[i].alignment = MessageAlignmentTop;
        }
        else if (strcmp(alignment, "bottom") == 0)
        {
            state->items[i].alignment = MessageAlignmentBottom;
        }
        else if (strcmp(alignment, "middle") == 0)
        {
            state->items[i].alignment = MessageAlignmentMiddle;
        }
        else
        {
            char buff[1024];
            snprintf(buff, sizeof(buff), "Invalid alignment provided for item %zu", i);
            log_error(buff);
            json_value_free(root_value);
            return NULL;
        }
    }

    state->item_count = item_count;
    state->selected = 0;

    if (json_object_has_value(root_object, "selected"))
    {
        state->selected = json_object_get_number(root_object, "selected");
        if (state->selected < 0)
        {
            state->selected = 0;
        }
        else if (state->selected >= item_count)
        {
            state->selected = item_count - 1;
        }
    }

    json_value_free(root_value);
    return state;
}

// handle_input interprets input events and mutates app state
void handle_input(struct AppState *state)
{
    if (!state->items_state->items[state->items_state->selected].image_exists && state->items_state->items[state->items_state->selected].background_image != NULL)
    {
        if (access(state->items_state->items[state->items_state->selected].background_image, F_OK) != -1)
        {
            state->items_state->items[state->items_state->selected].image_exists = true;
            state->redraw = 1;
        }
    }

    if (state->timeout_seconds < 0)
    {
        return;
    }

    if (increment_item_list_index)
    {
        pthread_mutex_lock(&increment_item_list_index_lock);
        increment_item_list_index = 0;
        state->items_state->selected += 1;
        bool should_return = false;
        if (state->items_state->selected >= state->items_state->item_count)
        {
            if (state->quit_after_last_item)
            {
                state->redraw = 0;
                state->quitting = 1;
                state->exit_code = ExitCodeSuccess;
                should_return = true;
            }
            else
            {
                state->items_state->selected = 0;
            }
        }
        state->redraw = 1;
        pthread_mutex_unlock(&increment_item_list_index_lock);
        if (should_return)
        {
            return;
        }
    }

    PAD_poll();

    bool is_action_button_pressed = false;
    bool is_confirm_button_pressed = false;
    bool is_cancel_button_pressed = false;
    bool is_inaction_button_pressed = false;
    if (PAD_justReleased(BTN_A))
    {
        if (strcmp(state->action_button, "A") == 0)
        {
            is_action_button_pressed = true;
        }
        else if (strcmp(state->confirm_button, "A") == 0)
        {
            is_confirm_button_pressed = true;
        }
        else if (strcmp(state->cancel_button, "A") == 0)
        {
            is_cancel_button_pressed = true;
        }
        else if (strcmp(state->inaction_button, "A") == 0)
        {
            is_inaction_button_pressed = true;
        }
    }
    else if (PAD_justReleased(BTN_B))
    {
        if (strcmp(state->action_button, "B") == 0)
        {
            is_action_button_pressed = true;
        }
        else if (strcmp(state->confirm_button, "B") == 0)
        {
            is_confirm_button_pressed = true;
        }
        else if (strcmp(state->cancel_button, "B") == 0)
        {
            is_cancel_button_pressed = true;
        }
        else if (strcmp(state->inaction_button, "B") == 0)
        {
            is_inaction_button_pressed = true;
        }
    }
    else if (PAD_justReleased(BTN_X))
    {
        if (strcmp(state->action_button, "X") == 0)
        {
            is_action_button_pressed = true;
        }
        else if (strcmp(state->confirm_button, "X") == 0)
        {
            is_confirm_button_pressed = true;
        }
        else if (strcmp(state->cancel_button, "X") == 0)
        {
            is_cancel_button_pressed = true;
        }
        else if (strcmp(state->inaction_button, "X") == 0)
        {
            is_inaction_button_pressed = true;
        }
    }
    else if (PAD_justReleased(BTN_Y))
    {
        if (strcmp(state->action_button, "Y") == 0)
        {
            is_action_button_pressed = true;
        }
        else if (strcmp(state->confirm_button, "Y") == 0)
        {
            is_confirm_button_pressed = true;
        }
        else if (strcmp(state->cancel_button, "Y") == 0)
        {
            is_cancel_button_pressed = true;
        }
        else if (strcmp(state->inaction_button, "Y") == 0)
        {
            is_inaction_button_pressed = true;
        }
    }

    if (is_action_button_pressed && state->action_show)
    {
        state->redraw = 0;
        state->quitting = 1;
        state->exit_code = ExitCodeActionButton;
        return;
    }

    if (is_confirm_button_pressed && state->confirm_show)
    {
        state->redraw = 0;
        state->quitting = 1;
        state->exit_code = ExitCodeSuccess;
        return;
    }

    if (is_cancel_button_pressed)
    {
        state->redraw = 0;
        state->quitting = 1;
        state->exit_code = ExitCodeCancelButton;
        return;
    }

    if (is_inaction_button_pressed && state->inaction_show)
    {
        state->redraw = 0;
        state->quitting = 1;
        state->exit_code = ExitCodeInactionButton;
        return;
    }

    if (PAD_justReleased(BTN_START))
    {
        state->redraw = 0;
        state->quitting = 1;
        state->exit_code = ExitCodeStartButton;
        return;
    }

    if (PAD_justRepeated(BTN_LEFT))
    {
        if (state->items_state->selected == 0 && !PAD_justPressed(BTN_LEFT))
        {
            state->redraw = 0;
        }
        else
        {
            state->items_state->selected -= 1;
            if (state->items_state->selected < 0)
            {
                state->items_state->selected = state->items_state->item_count - 1;
            }
            state->redraw = 1;
        }
    }
    else if (PAD_justRepeated(BTN_RIGHT))
    {
        if (state->items_state->selected == state->items_state->item_count - 1 && !PAD_justPressed(BTN_RIGHT))
        {
            state->redraw = 0;
        }
        else
        {
            state->items_state->selected += 1;
            if (state->items_state->selected >= state->items_state->item_count)
            {
                if (state->quit_after_last_item)
                {
                    state->redraw = 0;
                    state->quitting = 1;
                    state->exit_code = ExitCodeSuccess;
                    return;
                }
                state->items_state->selected = 0;
            }
            state->redraw = 1;
        }
    }
}

// turns a hex color (e.g. #000000) into an SDL_Color
SDL_Color hex_to_sdl_color(const char *hex)
{
    SDL_Color color = {0, 0, 0, 255};

    // Skip # if present
    if (hex[0] == '#')
    {
        hex++;
    }

    // Parse RGB values from hex string
    int r, g, b;
    if (sscanf(hex, "%02x%02x%02x", &r, &g, &b) == 3)
    {
        color.r = r;
        color.g = g;
        color.b = b;
    }

    return color;
}

// scale_surface manually scales a surface to a new width and height for SDL1
SDL_Surface *scale_surface(SDL_Surface *surface,
                           Uint16 width, Uint16 height)
{
    SDL_Surface *scaled = SDL_CreateRGBSurface(surface->flags,
                                               width,
                                               height,
                                               surface->format->BitsPerPixel,
                                               surface->format->Rmask,
                                               surface->format->Gmask,
                                               surface->format->Bmask,
                                               surface->format->Amask);

    int bpp = surface->format->BytesPerPixel;
    int *v = (int *)malloc(bpp * sizeof(int));

    for (int x = 0; x < width; x++)
    {
        for (int y = 0; y < height; y++)
        {
            int xo1 = x * surface->w / width;
            int xo2 = MAX((x + 1) * surface->w / width, xo1 + 1);
            int yo1 = y * surface->h / height;
            int yo2 = MAX((y + 1) * surface->h / height, yo1 + 1);
            int n = (xo2 - xo1) * (yo2 - yo1);

            for (int i = 0; i < bpp; i++)
                v[i] = 0;

            for (int xo = xo1; xo < xo2; xo++)
                for (int yo = yo1; yo < yo2; yo++)
                {
                    Uint8 *ps =
                        (Uint8 *)surface->pixels + yo * surface->pitch + xo * bpp;
                    for (int i = 0; i < bpp; i++)
                        v[i] += ps[i];
                }

            Uint8 *pd = (Uint8 *)scaled->pixels + y * scaled->pitch + x * bpp;
            for (int i = 0; i < bpp; i++)
                pd[i] = v[i] / n;
        }
    }

    free(v);

    return scaled;
}

// draw_screen interprets the app state and draws it to the screen
void draw_screen(SDL_Surface *screen, struct AppState *state)
{
    // render a background color
    char hex_color[1024] = "#000000";
    if (state->items_state->items[state->items_state->selected].background_color != NULL)
    {
        strncpy(hex_color, state->items_state->items[state->items_state->selected].background_color, sizeof(hex_color));
    }

    SDL_Color background_color = hex_to_sdl_color(hex_color);
    uint32_t color = SDL_MapRGBA(screen->format, background_color.r, background_color.g, background_color.b, 255);
    SDL_FillRect(screen, NULL, color);

    // check if there is an image and it is accessible
    if (state->items_state->items[state->items_state->selected].background_image != NULL)
    {
        SDL_Surface *surface = IMG_Load(state->items_state->items[state->items_state->selected].background_image);
        if (surface)
        {
            int imgW = surface->w, imgH = surface->h;

            // Compute scale factor
            float scaleX = (float)(FIXED_WIDTH - 2 * PADDING) / imgW;
            float scaleY = (float)(FIXED_HEIGHT - 2 * PADDING) / imgH;
            float scale = (scaleX < scaleY) ? scaleX : scaleY;

            // Ensure upscaling only when the image is smaller than the screen
            if (imgW * scale < FIXED_WIDTH - 2 * PADDING && imgH * scale < FIXED_HEIGHT - 2 * PADDING)
            {
                scale = (scaleX > scaleY) ? scaleX : scaleY;
            }

            // Compute target dimensions
            int dstW = imgW * scale;
            int dstH = imgH * scale;

            int dstX = (FIXED_WIDTH - dstW) / 2;
            int dstY = (FIXED_HEIGHT - dstH) / 2;
            if (imgW == FIXED_WIDTH && imgH == FIXED_HEIGHT)
            {
                dstW = FIXED_WIDTH;
                dstH = FIXED_HEIGHT;
                dstX = 0;
                dstY = 0;
            }

            // Compute destination rectangle
            SDL_Rect dstRect = {dstX, dstY, dstW, dstH};
#ifdef USE_SDL2
            SDL_BlitScaled(surface, NULL, screen, &dstRect);
#else
            if (imgW == FIXED_WIDTH && imgH == FIXED_HEIGHT)
            {
                SDL_BlitSurface(surface, NULL, screen, &dstRect);
            }
            else
            {
                SDL_Surface *scaled = scale_surface(surface, dstW, dstH);
                SDL_BlitSurface(scaled, NULL, screen, &dstRect);
                SDL_FreeSurface(scaled);
            }
#endif
            SDL_FreeSurface(surface);
        }
    }

    // draw the button group on the button-right
    // only two buttons can be displayed at a time
    if (state->confirm_show && strcmp(state->confirm_button, "") != 0)
    {
        if (state->cancel_show && strcmp(state->cancel_button, "") != 0)
        {
            GFX_blitButtonGroup((char *[]){state->cancel_button, state->cancel_text, state->confirm_button, state->confirm_text, NULL}, 1, screen, 1);
        }
        else
        {
            GFX_blitButtonGroup((char *[]){state->confirm_button, state->confirm_text, NULL}, 1, screen, 1);
        }
    }
    else if (state->cancel_show)
    {
        GFX_blitButtonGroup((char *[]){state->cancel_button, state->cancel_text, NULL}, 1, screen, 1);
    }

    int initial_padding = 0;
    if (state->show_time_left && state->timeout_seconds > 0)
    {
        struct timeval current_time;
        gettimeofday(&current_time, NULL);

        int time_left = state->timeout_seconds - (current_time.tv_sec - state->start_time.tv_sec);
        if (time_left <= 0)
        {
            time_left = 0;
        }

        char time_left_str[1024];
        if (time_left == 1)
        {
            snprintf(time_left_str, sizeof(time_left_str), "Time left: %d second", time_left);
        }
        else
        {
            snprintf(time_left_str, sizeof(time_left_str), "Time left: %d seconds", time_left);
        }

        SDL_Surface *text = TTF_RenderUTF8_Blended(state->fonts.small, time_left_str, COLOR_WHITE);
        SDL_Rect pos = {
            SCALE1(PADDING),
            SCALE1(PADDING),
            text->w,
            text->h};
        SDL_BlitSurface(text, NULL, screen, &pos);

        initial_padding = text->h + SCALE1(PADDING);
    }

    int message_padding = SCALE1(PADDING + BUTTON_PADDING);

    // get the width and height of every word in the message
    struct Message words[1024];
    int word_count = 0;
    char original_message[1024];
    strncpy(original_message, state->items_state->items[state->items_state->selected].text, sizeof(original_message));
    char *word = strtok(original_message, " ");
    int word_height = 0;
    while (word != NULL)
    {
        int word_width;
        strtrim(word);
        if (strcmp(word, "") == 0)
        {
            continue;
        }

        TTF_SizeUTF8(state->fonts.large, word, &word_width, &word_height);
        strncpy(words[word_count].message, word, sizeof(words[word_count].message));
        words[word_count].width = word_width;
        word_count++;
        word = strtok(NULL, " ");
    }

    int letter_width = 0;
    TTF_SizeUTF8(state->fonts.large, "A", &letter_width, NULL);

    // construct a list of messages that can be displayed on a single line
    // if the message is too long to be displayed on a single line,
    // the message will be wrapped onto multiple lines
    struct Message messages[MAIN_ROW_COUNT];
    for (int i = 0; i < MAIN_ROW_COUNT; i++)
    {
        strncpy(messages[i].message, "", sizeof(messages[i].message));
        messages[i].width = 0;
    }

    int message_count = 0;
    int current_message_index = 0;
    for (int i = 0; i < word_count; i++)
    {
        if (current_message_index >= MAIN_ROW_COUNT)
        {
            break;
        }

        int potential_width = messages[current_message_index].width + words[i].width;
        if (i > 0)
        {
            potential_width += letter_width;
        }

        if (messages[current_message_index].width == 0)
        {
            strncpy(messages[current_message_index].message, words[i].message, sizeof(messages[current_message_index].message));
            messages[current_message_index].width = words[i].width;
        }
        else if (potential_width <= FIXED_WIDTH - 2 * message_padding)
        {
            if (messages[current_message_index].width == 0)
            {
                strncpy(messages[current_message_index].message, words[i].message, sizeof(messages[current_message_index].message));
            }
            else
            {
                char messageBuf[256];
                snprintf(messageBuf, sizeof(messageBuf), "%s %s", messages[current_message_index].message, words[i].message);

                strncpy(messages[current_message_index].message, messageBuf, sizeof(messages[current_message_index].message));
            }
            messages[current_message_index].width += words[i].width;
        }
        else
        {
            current_message_index++;
            message_count++;
            strncpy(messages[current_message_index].message, words[i].message, sizeof(messages[current_message_index].message));
            messages[current_message_index].width = words[i].width;
        }
    }

    int messages_height = (current_message_index + 1) * word_height + (SCALE1(PADDING) * current_message_index);
    // default to the middle of the screen
    int current_message_y = (screen->h - messages_height) / 2;
    if (state->items_state->items[state->items_state->selected].alignment == MessageAlignmentTop)
    {
        current_message_y = SCALE1(PADDING) + initial_padding;
    }
    else if (state->items_state->items[state->items_state->selected].alignment == MessageAlignmentBottom)
    {
        current_message_y = screen->h - messages_height - SCALE1(PADDING) - initial_padding;
    }

    for (int i = 0; i <= message_count; i++)
    {
        char *message = messages[i].message;
        if (message == NULL)
        {
            continue;
        }

        int width = messages[i].width;
        SDL_Surface *text = TTF_RenderUTF8_Blended(state->fonts.large, message, COLOR_WHITE);
        if (text == NULL)
        {
            continue;
        }

        SDL_Rect pos = {
            ((screen->w - text->w) / 2),
            current_message_y + PADDING,
            text->w,
            text->h};

        if (state->items_state->items[state->items_state->selected].show_pill)
        {
            SDL_Rect pill_rect = {
                pos.x - SCALE1(PADDING * 2),
                pos.y - SCALE1(PADDING),
                text->w + SCALE1(PADDING * 4),
                SCALE1(PILL_SIZE)};
            GFX_blitPill(ASSET_BLACK_PILL, screen, &pill_rect);
        }

        SDL_BlitSurface(text, NULL, screen, &pos);
        current_message_y += word_height + SCALE1(PADDING);
    }
    if (state->action_show && strcmp(state->action_button, "") != 0)
    {
        if (state->inaction_show && strcmp(state->inaction_button, "") != 0)
        {
            GFX_blitButtonGroup((char *[]){state->inaction_button, state->inaction_text, state->action_button, state->action_text, NULL}, 0, screen, 0);
        }
        else
        {
            GFX_blitButtonGroup((char *[]){state->action_button, state->action_text, NULL}, 0, screen, 0);
        }
    }
    else if (state->inaction_show && strcmp(state->inaction_button, "") != 0)
    {
        GFX_blitButtonGroup((char *[]){state->inaction_button, state->inaction_text, NULL}, 0, screen, 0);
    }

    // don't forget to reset the should_redraw flag
    state->redraw = 0;
}

bool open_fonts(struct AppState *state)
{
    if (state->fonts.font_path == NULL)
    {
        log_error("No font path provided");
        return false;
    }

    // check if the font path is valid
    if (access(state->fonts.font_path, F_OK) == -1)
    {
        log_error("Invalid font path provided");
        return false;
    }

    state->fonts.large = TTF_OpenFont(state->fonts.font_path, SCALE1(state->fonts.size));
    if (state->fonts.large == NULL)
    {
        char buff[1024];
        snprintf(buff, sizeof(buff), "Failed to open large font: %s", TTF_GetError());
        log_error(buff);
        return false;
    }
    TTF_SetFontStyle(state->fonts.large, TTF_STYLE_BOLD);

    state->fonts.small = TTF_OpenFont(state->fonts.font_path, SCALE1(FONT_SMALL));
    if (state->fonts.small == NULL)
    {
        char buff[1024];
        snprintf(buff, sizeof(buff), "Failed to open small font: %s", TTF_GetError());
        log_error(buff);
        return false;
    }

    return true;
}

void signal_handler(int signal)
{
    // if the signal is a ctrl+c, exit with code 130
    if (signal == SIGINT)
    {
        exit(ExitCodeKeyboardInterrupt);
    }
    else if (signal == SIGTERM)
    {
        exit(ExitCodeSigterm);
    }
    else if (signal == SIGUSR1)
    {
        increment_item_list_index = 1;
    }
    else
    {
        exit(ExitCodeError);
    }
}

// parse_arguments parses the arguments using getopt and updates the app state
// supports the following flags:
// - --action-button <button> (default: "")
// - --action-text <text> (default: "ACTION")
// - --action-show (default: false)
// - --background-image <path> (default: empty string)
// - --background-color <hex> (default: empty string)
// - --confirm-button <button> (default: "A")
// - --confirm-text <text> (default: "SELECT")
// - --confirm-show (default: false)
// - --cancel-button <button> (default: "B")
// - --cancel-text <text> (default: "BACK")
// - --cancel-show (default: false)
// - --inaction-button <button> (default: empty string)
// - --inaction-text <text> (default: "OTHER")
// - --inaction-show (default: false)
// - --file <path> (default: empty string)
// - --item-key <key> (default: "items")
// - --message <message> (default: empty string)
// - --message-alignment <alignment> (default: middle)
// - --font <path> (default: empty string)
// - --font-size <size> (default: FONT_LARGE)
// - --quit-after-last-item (default: false)
// - --show-hardware-group (default: false)
// - --show-pill (default: false)
// - --show-time-left (default: false)
// - --timeout <seconds> (default: 1)
bool parse_arguments(struct AppState *state, int argc, char *argv[])
{
    static struct option long_options[] = {
        {"action-button", required_argument, 0, 'a'},
        {"action-text", required_argument, 0, 'A'},
        {"background-image", required_argument, 0, 'b'},
        {"background-color", required_argument, 0, 'B'},
        {"confirm-button", required_argument, 0, 'c'},
        {"confirm-text", required_argument, 0, 'C'},
        {"cancel-button", required_argument, 0, 'd'},
        {"cancel-text", required_argument, 0, 'D'},
        {"inaction-button", required_argument, 0, 'i'},
        {"inaction-text", required_argument, 0, 'I'},
        {"file", required_argument, 0, 'E'},
        {"font-default", required_argument, 0, 'f'},
        {"font-size-default", required_argument, 0, 'F'},
        {"item-key", required_argument, 0, 'K'},
        {"message", required_argument, 0, 'm'},
        {"message-alignment", required_argument, 0, 'M'},
        {"quit-after-last-item", no_argument, 0, 'Q'},
        {"show-pill", no_argument, 0, 'P'},
        {"show-hardware-group", no_argument, 0, 'S'},
        {"show-time-left", no_argument, 0, 'T'},
        {"timeout", required_argument, 0, 't'},
        {"confirm-show", no_argument, 0, 'W'},
        {"cancel-show", no_argument, 0, 'X'},
        {"action-show", no_argument, 0, 'Y'},
        {"inaction-show", no_argument, 0, 'Z'},
        {0, 0, 0, 0}};

    int opt;
    char *font_path = NULL;
    char message[1024];
    char alignment[1024];
    while ((opt = getopt_long(argc, argv, "a:A:b:B:c:C:d:D:E:f:F:i:I:K:m:M::t:QPSTWYXZ", long_options, NULL)) != -1)
    {
        switch (opt)
        {
        case 'a':
            strncpy(state->action_button, optarg, sizeof(state->action_button));
            break;
        case 'A':
            strncpy(state->action_text, optarg, sizeof(state->action_text));
            break;
        case 'b':
            strncpy(state->background_image, optarg, sizeof(state->background_image));
            break;
        case 'B':
            strncpy(state->background_color, optarg, sizeof(state->background_color));
            break;
        case 'c':
            strncpy(state->confirm_button, optarg, sizeof(state->confirm_button));
            break;
        case 'C':
            strncpy(state->confirm_text, optarg, sizeof(state->confirm_text));
            break;
        case 'd':
            strncpy(state->cancel_button, optarg, sizeof(state->cancel_button));
            break;
        case 'D':
            strncpy(state->cancel_text, optarg, sizeof(state->cancel_text));
            break;
        case 'E':
            strncpy(state->file, optarg, sizeof(state->file));
            break;
        case 'f':
            font_path = optarg;
            break;
        case 'F':
            state->fonts.size = atoi(optarg);
            break;
        case 'i':
            strncpy(state->inaction_button, optarg, sizeof(state->inaction_button));
            break;
        case 'I':
            strncpy(state->inaction_text, optarg, sizeof(state->inaction_text));
            break;
        case 'K':
            strncpy(state->item_key, optarg, sizeof(state->item_key));
            break;
        case 'm':
            strncpy(message, optarg, sizeof(message));
            break;
        case 'M':
            strncpy(alignment, optarg, sizeof(alignment));
            break;
        case 'Q':
            state->quit_after_last_item = true;
            break;
        case 'P':
            state->show_pill = true;
            break;
        case 'S':
            state->show_hardware_group = true;
            break;
        case 't':
            state->timeout_seconds = atoi(optarg);
            break;
        case 'T':
            state->show_time_left = true;
            break;
        case 'W':
            state->confirm_show = true;
            break;
        case 'X':
            state->cancel_show = true;
            break;
        case 'Y':
            state->action_show = true;
            break;
        case 'Z':
            state->inaction_show = true;
            break;
        default:
            return false;
        }
    }

    enum MessageAlignment default_alignment = MessageAlignmentMiddle;
    if (strcmp(alignment, "top") == 0)
    {
        default_alignment = MessageAlignmentTop;
    }
    else if (strcmp(alignment, "bottom") == 0)
    {
        default_alignment = MessageAlignmentBottom;
    }
    else if (strcmp(alignment, "middle") == 0 || strcmp(alignment, "") == 0)
    {
        default_alignment = MessageAlignmentMiddle;
    }
    else
    {
        log_error("Invalid message alignment provided");
        return false;
    }

    if (strlen(message) > 0)
    {
        struct ItemsState *items_state = malloc(sizeof(struct ItemsState));
        items_state->items = malloc(sizeof(struct Item) * 1);
        items_state->items[0].text = strdup(message);
        items_state->items[0].background_color = "#000000";
        items_state->items[0].background_image = NULL;
        items_state->items[0].image_exists = false;
        items_state->items[0].show_pill = state->show_pill;

        if (strcmp(state->background_color, "") != 0)
        {
            items_state->items[0].background_color = strdup(state->background_color);
        }

        if (strcmp(state->background_image, "") != 0)
        {
            items_state->items[0].background_image = strdup(state->background_image);
            items_state->items[0].image_exists = access(state->background_image, F_OK) != -1;
        }

        items_state->items[0].alignment = default_alignment;
        items_state->item_count = 1;
        items_state->selected = 0;
        state->items_state = items_state;
    }
    else if (strcmp(state->file, "") != 0)
    {
        state->items_state = ItemsState_New(state->file, state->item_key, state->background_image, state->background_color, state->show_pill, default_alignment);
        if (state->items_state == NULL)
        {
            log_error("Failed to hydrate display states");
            return false;
        }
    }
    else
    {
        log_error("No message or file provided");
        return false;
    }

    if (font_path != NULL)
    {
        // check if the font path is valid
        if (access(font_path, F_OK) == -1)
        {
            log_error("Invalid font path provided");
            return false;
        }

        state->fonts.font_path = strdup(font_path);
    }
    else
    {
        state->fonts.font_path = strdup(FONT_PATH);
    }

    // Apply default values for certain buttons and texts
    if (strcmp(state->action_button, "") == 0)
    {
        strncpy(state->action_button, "", sizeof(state->action_button));
    }

    if (strcmp(state->action_text, "") == 0)
    {
        strncpy(state->action_text, "ACTION", sizeof(state->action_text));
    }

    if (strcmp(state->cancel_button, "") == 0)
    {
        strncpy(state->cancel_button, "B", sizeof(state->cancel_button));
    }

    if (strcmp(state->confirm_text, "") == 0)
    {
        strncpy(state->confirm_text, "SELECT", sizeof(state->confirm_text));
    }

    if (strcmp(state->cancel_text, "") == 0)
    {
        strncpy(state->cancel_text, "BACK", sizeof(state->cancel_text));
    }

    if (strcmp(state->inaction_text, "") == 0)
    {
        strncpy(state->inaction_text, "OTHER", sizeof(state->inaction_text));
    }

    // validate that hardware buttons aren't assigned to more than once
    bool a_button_assigned = false;
    bool b_button_assigned = false;
    bool x_button_assigned = false;
    bool y_button_assigned = false;
    bool i_button_assigned = false;
    if (strcmp(state->action_button, "A") == 0)
    {
        a_button_assigned = true;
    }
    if (strcmp(state->cancel_button, "A") == 0)
    {
        if (a_button_assigned)
        {
            log_error("A button cannot be assigned to more than one button");
            return false;
        }

        a_button_assigned = true;
    }
    if (strcmp(state->confirm_button, "A") == 0)
    {
        if (a_button_assigned)
        {
            log_error("A button cannot be assigned to more than one button");
            return false;
        }

        a_button_assigned = true;
    }
    if (strcmp(state->inaction_button, "A") == 0)
    {
        if (a_button_assigned)
        {
            log_error("A button cannot be assigned to more than one button");
            return false;
        }

        a_button_assigned = true;
    }

    if (strcmp(state->action_button, "B") == 0)
    {
        b_button_assigned = true;
    }
    if (strcmp(state->cancel_button, "B") == 0)
    {
        if (b_button_assigned)
        {
            log_error("B button cannot be assigned to more than one button");
            return false;
        }

        b_button_assigned = true;
    }
    if (strcmp(state->confirm_button, "B") == 0)
    {
        if (b_button_assigned)
        {
            log_error("B button cannot be assigned to more than one button");
            return false;
        }

        b_button_assigned = true;
    }
    if (strcmp(state->inaction_button, "B") == 0)
    {
        if (b_button_assigned)
        {
            log_error("B button cannot be assigned to more than one button");
            return false;
        }

        b_button_assigned = true;
    }

    if (strcmp(state->action_button, "X") == 0)
    {
        x_button_assigned = true;
    }
    if (strcmp(state->cancel_button, "X") == 0)
    {
        if (x_button_assigned)
        {
            log_error("X button cannot be assigned to more than one button");
            return false;
        }

        x_button_assigned = true;
    }
    if (strcmp(state->confirm_button, "X") == 0)
    {
        if (x_button_assigned)
        {
            log_error("X button cannot be assigned to more than one button");
            return false;
        }

        x_button_assigned = true;
    }
    if (strcmp(state->inaction_button, "X") == 0)
    {
        if (x_button_assigned)
        {
            log_error("X button cannot be assigned to more than one button");
            return false;
        }

        x_button_assigned = true;
    }

    if (strcmp(state->action_button, "Y") == 0)
    {
        y_button_assigned = true;
    }
    if (strcmp(state->cancel_button, "Y") == 0)
    {
        if (y_button_assigned)
        {
            log_error("Y button cannot be assigned to more than one button");
            return false;
        }

        y_button_assigned = true;
    }
    if (strcmp(state->confirm_button, "Y") == 0)
    {
        if (y_button_assigned)
        {
            log_error("Y button cannot be assigned to more than one button");
            return false;
        }

        y_button_assigned = true;
    }
    if (strcmp(state->inaction_button, "Y") == 0)
    {
        if (y_button_assigned)
        {
            log_error("Y button cannot be assigned to more than one button");
            return false;
        }

        y_button_assigned = true;
    }

    // validate that the confirm and cancel buttons are valid
    if (strcmp(state->confirm_button, "A") != 0 && strcmp(state->confirm_button, "B") != 0 && strcmp(state->confirm_button, "X") != 0 && strcmp(state->confirm_button, "Y") != 0 && strcmp(state->confirm_button, "") != 0)
    {
        log_error("Invalid confirm button provided");
        return false;
    }
    if (strcmp(state->cancel_button, "A") != 0 && strcmp(state->cancel_button, "B") != 0 && strcmp(state->cancel_button, "X") != 0 && strcmp(state->cancel_button, "Y") != 0 && strcmp(state->cancel_button, "") != 0)
    {
        log_error("Invalid cancel button provided");
        return false;
    }
    if (strcmp(state->action_button, "A") != 0 && strcmp(state->action_button, "B") != 0 && strcmp(state->action_button, "X") != 0 && strcmp(state->action_button, "Y") != 0 && strcmp(state->action_button, "") != 0)
    {
        log_error("Invalid action button provided");
        return false;
    }
    if (strcmp(state->inaction_button, "A") != 0 && strcmp(state->inaction_button, "B") != 0 && strcmp(state->inaction_button, "X") != 0 && strcmp(state->inaction_button, "Y") != 0 && strcmp(state->inaction_button, "") != 0)
    {
        log_error("Invalid inaction button provided");
        return false;
    }

    return true;
}

// suppress_output suppresses stdout and stderr
// returns a single integer containing both file descriptors
int suppress_output(void)
{
    int stdout_fd = dup(STDOUT_FILENO);
    int stderr_fd = dup(STDERR_FILENO);

    int dev_null_fd = open("/dev/null", O_WRONLY);
    dup2(dev_null_fd, STDOUT_FILENO);
    dup2(dev_null_fd, STDERR_FILENO);
    close(dev_null_fd);

    return (stdout_fd << 16) | stderr_fd;
}

// restore_output restores stdout and stderr to the original file descriptors
void restore_output(int saved_fds)
{
    int stdout_fd = (saved_fds >> 16) & 0xFFFF;
    int stderr_fd = saved_fds & 0xFFFF;

    fflush(stdout);
    fflush(stderr);

    dup2(stdout_fd, STDOUT_FILENO);
    dup2(stderr_fd, STDERR_FILENO);

    close(stdout_fd);
    close(stderr_fd);
}

// swallow_stdout_from_function swallows stdout from a function
// this is useful for suppressing output from a function
// that we don't want to see in the log file
// the InitSettings() function is an example of this (some implementations print to stdout)
void swallow_stdout_from_function(void (*func)(void))
{
    int saved_fds = suppress_output();

    func();

    restore_output(saved_fds);
}

// init initializes the app state
// everything is placed here as MinUI sometimes logs to stdout
// and the logging happens depending on the platform
void init()
{
    // set the cpu speed to the menu speed
    // this is done here to ensure we downclock
    // the menu (no need to draw power unnecessarily)
    PWR_setCPUSpeed(CPU_SPEED_MENU);

    // initialize:
    // - the screen, allowing us to draw to it
    // - input from the pad/joystick/buttons/etc.
    // - power management
    // - sync hardware settings (brightness, hdmi, speaker, etc.)
    if (screen == NULL)
    {
        screen = GFX_init(MODE_MAIN);
    }
    PAD_init();
    PWR_init();
    InitSettings();
}

// destruct cleans up the app state in reverse order
void destruct()
{
    QuitSettings();
    PWR_quit();
    PAD_quit();
    GFX_quit();
}

// main is the entry point for the app
int main(int argc, char *argv[])
{
    // Initialize app state
    char default_action_button[1024] = "";
    char default_action_text[1024] = "ACTION";
    char default_background_image[1024] = "";
    char default_background_color[1024] = "#000000";
    char default_cancel_button[1024] = "B";
    char default_cancel_text[1024] = "BACK";
    char default_confirm_button[1024] = "A";
    char default_confirm_text[1024] = "SELECT";
    char default_inaction_button[1024] = "";
    char default_inaction_text[1024] = "OTHER";
    char default_file[1024] = "";
    char default_item_key[1024] = "items";
    char default_message[1024] = "";
    struct AppState state = {
        .redraw = 1,
        .quitting = 0,
        .exit_code = ExitCodeSuccess,
        .show_hardware_group = false,
        .timeout_seconds = 0,
        .fonts = {
            .size = FONT_LARGE,
            .large = NULL,
            .medium = NULL,
            .font_path = NULL,
        },
        .action_show = false,
        .confirm_show = false,
        .cancel_show = false,
        .inaction_show = false,
        .quit_after_last_item = false,
        .show_time_left = false,
        .items_state = NULL,
        .start_time = 0,
        .show_pill = false,
    };

    // assign the default values to the app state
    strncpy(state.action_button, default_action_button, sizeof(state.action_button));
    strncpy(state.action_text, default_action_text, sizeof(state.action_text));
    strncpy(state.background_image, default_background_image, sizeof(state.background_image));
    strncpy(state.background_color, default_background_color, sizeof(state.background_color));
    strncpy(state.cancel_button, default_cancel_button, sizeof(state.cancel_button));
    strncpy(state.cancel_text, default_cancel_text, sizeof(state.cancel_text));
    strncpy(state.confirm_button, default_confirm_button, sizeof(state.confirm_button));
    strncpy(state.confirm_text, default_confirm_text, sizeof(state.confirm_text));
    strncpy(state.inaction_button, default_inaction_button, sizeof(state.inaction_button));
    strncpy(state.inaction_text, default_inaction_text, sizeof(state.inaction_text));
    strncpy(state.file, default_file, sizeof(state.file));
    strncpy(state.item_key, default_item_key, sizeof(state.item_key));

    // parse the arguments
    if (!parse_arguments(&state, argc, argv))
    {
        return ExitCodeError;
    }

    swallow_stdout_from_function(init);

    struct sigaction sa = {
        .sa_handler = signal_handler,
        .sa_flags = SA_RESTART};
    sigemptyset(&sa.sa_mask);
    sigaction(SIGINT, &sa, NULL);
    sigaction(SIGTERM, &sa, NULL);
    sigaction(SIGUSR1, &sa, NULL);

    if (!open_fonts(&state))
    {
        return ExitCodeError;
    }

    // get initial wifi state
    int was_online = PLAT_isOnline();

    // get the current time
    gettimeofday(&state.start_time, NULL);

    int show_setting = 0; // 1=brightness,2=volume

    if (state.timeout_seconds <= 0)
    {
        PWR_disableAutosleep();
    }

    while (!state.quitting)
    {
        // start the frame to ensure GFX_sync() works
        // on devices that don't support vsync
        GFX_startFrame();

        // handle turning the on/off screen on/off
        // as well as general power management
        PWR_update(&state.redraw, &show_setting, NULL, NULL);

        // check if the device is on wifi
        // redraw if the wifi state changed
        // and then update our state
        int is_online = PLAT_isOnline();
        if (was_online != is_online)
        {
            state.redraw = 1;
        }
        was_online = is_online;

        // handle any input events
        handle_input(&state);

        // redraw the screen if there has been a change
        if (state.redraw)
        {
            // clear the screen at the beginning of each loop
            GFX_clear(screen);

            if (state.show_hardware_group)
            {
                // draw the hardware information in the top-right
                GFX_blitHardwareGroup(screen, show_setting);
                // draw the setting hints
                if (show_setting && !GetHDMI())
                {
                    GFX_blitHardwareHints(screen, show_setting);
                }
                else
                {
                    GFX_blitButtonGroup((char *[]){BTN_SLEEP == BTN_POWER ? "POWER" : "MENU", "SLEEP", NULL}, 0, screen, 0);
                }
            }

            // your draw logic goes here
            draw_screen(screen, &state);

            // Takes the screen buffer and displays it on the screen
            GFX_flip(screen);
        }
        else
        {
            // Slows down the frame rate to match the refresh rate of the screen
            // when the screen is not being redrawn
            GFX_sync();
        }

        // if the sleep seconds is larger than 0, check if the sleep has expired
        if (state.timeout_seconds > 0)
        {
            struct timeval current_time;
            gettimeofday(&current_time, NULL);
            if (current_time.tv_sec - state.start_time.tv_sec >= state.timeout_seconds)
            {
                state.exit_code = ExitCodeTimeout;
                state.quitting = 1;
            }

            if (current_time.tv_sec != state.start_time.tv_sec && state.show_time_left)
            {
                state.redraw = 1;
            }
        }
    }

    swallow_stdout_from_function(destruct);

    // exit the program
    return state.exit_code;
}
