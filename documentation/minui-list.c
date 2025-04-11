#include <fcntl.h>
#include <getopt.h>
#include <msettings.h>
#include <parson/parson.h>
#include <signal.h>
#include <stdbool.h>
#include <stdio.h>
#include <unistd.h>
#ifdef USE_SDL2
#include <SDL2/SDL_ttf.h>
#else
#include <SDL/SDL_ttf.h>
#endif

#include "defines.h"
#include "api.h"
#include "utils.h"

SDL_Surface *screen = NULL;

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

#define OPTION_PADDING 8

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

struct ListItemFeature
{
    // the background color to use for the list
    char background_color[1024];
    // path to the background image to use for the list
    char background_image[1024];
    // whether the background image exists
    bool background_image_exists;
    // whether the item can be disabled
    bool can_disable;
    // the confirm text to display on the confirm button
    char confirm_text[1024];
    // whether the item is disabled
    bool disabled;
    // whether to draw arrows around the item
    bool draw_arrows;
    // whether to hide the action button when the item is selected
    bool hide_action;
    // whether to hide the cancel button when the item is selected
    bool hide_cancel;
    // whether to hide the confirm button when the item is selected
    bool hide_confirm;
    // whether or not the item is a header
    bool is_header;
    // whether or not the item is unselectable
    bool unselectable;
    // alignment of the item text ('left', 'center', 'right')
    char alignment[1024];

    // whether the item has a background_color field
    bool has_background_color;
    // whether the item has a background_image field
    bool has_background_image;
    // whether the item has a can_disable field
    bool has_can_disable;
    // whether the item has a confirm_text field
    bool has_confirm_text;
    // whether the item has a disabled field
    bool has_disabled;
    // whether the item has a draw_arrows field
    bool has_draw_arrows;
    // whether the item has a hide_action field
    bool has_hide_action;
    // whether the item has a hide_cancel field
    bool has_hide_cancel;
    // whether the item has a hide_confirm field
    bool has_hide_confirm;
    // whether the item has a is_header field
    bool has_is_header;
    // whether the item has a unselectable field
    bool has_unselectable;
    // whether the item has a alignment field
    bool has_alignment;
};

// ListItem holds the configuration for a list item
struct ListItem
{
    // the name of the item
    char *name;
    // whether the item has features
    bool has_features;
    // whether the item has options field
    bool has_options;
    // whether the item has a selected field
    bool has_selected;
    // the number of options for the item
    int option_count;
    // a list of char options for the item
    char **options;
    // the selected option index
    int selected;
    // the initial selected option index
    int initial_selected;
    // the features of the item
    struct ListItemFeature features;
};

// ListState holds the state of the list
struct ListState
{
    // array of list items
    struct ListItem *items;
    // number of items in the list
    size_t item_count;
    // rendering state
    // index of first visible item
    int first_visible;
    // index of last visible item
    int last_visible;
    // index of currently selected item
    int selected;

    // whether or not any items in the list have options
    bool has_options;
};

// Fonts holds the fonts for the list
struct Fonts
{
    // the large font to use for the list
    TTF_Font *large;
    // the medium font to use for the list
    TTF_Font *medium;
    // the path to the default font to use for the list
    char *default_font;
    // the path to the large font to use for the list
    char *large_font;
    // the path to the medium font to use for the list
    char *medium_font;
};

// AppState holds the current state of the application
struct AppState
{
    // the exit code to return
    int exit_code;
    // whether the app should exit
    int quitting;
    // whether the screen needs to be redrawn
    int redraw;
    // whether to show the hardware group
    int show_hardware_group;
    // whether to show the brightness or hardware state
    int show_brightness_setting;
    // the button to display on the Action button
    char action_button[1024];
    // the text to display on the Action button
    char action_text[1024];
    // the background image to display
    char background_image[1024];
    // the background color to display
    char background_color[1024];
    // the button to display on the Confirm button
    char confirm_button[1024];
    // the text to display on the Confirm button
    char confirm_text[1024];
    // the button to display on the Cancel button
    char cancel_button[1024];
    // the text to display on the Cancel button
    char cancel_text[1024];
    // whether to disable sleep
    bool disable_auto_sleep;
    // the button to display on the Enable button
    char enable_button[1024];
    // the path to the JSON file
    char file[1024];
    // the format to read the input from
    char format[1024];
    // the key to the items array in the JSON file
    char item_key[1024];
    // the title of the list page
    char title[1024];
    // the title alignment ('left', 'center', 'right')
    char title_alignment[1024];
    // the location to write the value to
    char write_location[1024];
    // the value to write (selected, state)
    char write_value[1024];
    // the fonts to use for the list
    struct Fonts fonts;
    // the state of the list
    struct ListState *list_state;
};

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

char *read_file(const char *filename)
{
    FILE *file = fopen(filename, "r");
    if (!file)
    {
        return NULL;
    }

    // Get file size
    fseek(file, 0, SEEK_END);
    long file_size = ftell(file);
    fseek(file, 0, SEEK_SET);

    // Allocate buffer with extra byte for null terminator
    char *contents = malloc(file_size + 1);
    if (!contents)
    {
        fclose(file);
        return NULL;
    }

    // Read file contents
    size_t bytes_read = fread(contents, 1, file_size, file);
    fclose(file);

    if (bytes_read != file_size)
    {
        free(contents);
        return NULL;
    }

    // Add null terminator
    contents[file_size] = '\0';

    return contents;
}

// ListState_New creates a new ListState from a JSON file
struct ListState *ListState_New(const char *filename, const char *format, const char *item_key, const char *title, const char *confirm_text, const char *default_background_image, const char *default_background_color)
{
    struct ListState *state = malloc(sizeof(struct ListState));

    int max_row_count = MAIN_ROW_COUNT;
    if (strlen(title) > 0)
    {
        max_row_count -= 1;
    }

    if (strcmp(format, "text") == 0)
    {
        char *contents = NULL;
        if (strcmp(filename, "-") == 0)
        {
            contents = read_stdin();
        }
        else
        {
            contents = read_file(filename);
        }

        if (contents == NULL)
        {
            log_error("Failed to read file or stdin");
            free(state);
            return NULL;
        }

        // Count number of non-empty lines
        size_t item_count = 0;
        char *line_start = contents;
        while (*line_start != '\0')
        {
            char *line_end = strchr(line_start, '\n');
            if (!line_end)
            {
                line_end = line_start + strlen(line_start);
            }

            // Check if line has non-whitespace content
            char *p;
            for (p = line_start; p < line_end && isspace(*p); p++)
                ;
            if (p < line_end)
            {
                item_count++;
            }

            if (*line_end == '\0')
            {
                break;
            }
            line_start = line_end + 1;
        }

        // Allocate array for items
        state->items = malloc(sizeof(struct ListItem) * item_count);
        state->item_count = item_count;
        state->last_visible = (item_count < max_row_count) ? item_count : max_row_count;
        state->first_visible = 0;
        state->selected = 0;

        // Add non-empty lines to items array
        size_t item_index = 0;
        line_start = contents;
        while (*line_start != '\0')
        {
            char *line_end = strchr(line_start, '\n');
            if (!line_end)
            {
                line_end = line_start + strlen(line_start);
            }

            // Check if line has non-whitespace content
            char *p;
            for (p = line_start; p < line_end && isspace(*p); p++)
                ;
            if (p < line_end)
            {
                size_t line_len = line_end - line_start;
                char *line = malloc(line_len + 1);
                memcpy(line, line_start, line_len);
                line[line_len] = '\0';
                state->items[item_index].name = line;
                state->items[item_index].has_features = false;
                state->items[item_index].has_options = false;
                state->items[item_index].has_selected = false;
                state->items[item_index].option_count = 0;
                state->items[item_index].options = NULL;
                state->items[item_index].selected = 0;
                state->items[item_index].initial_selected = 0;
                state->items[item_index].features = (struct ListItemFeature){
                    .background_color = "",
                    .background_image = "",
                    .background_image_exists = false,
                    .can_disable = false,
                    .confirm_text = "",
                    .disabled = false,
                    .draw_arrows = false,
                    .hide_action = false,
                    .hide_cancel = false,
                    .hide_confirm = false,
                    .is_header = false,
                    .unselectable = false,
                    .alignment = "",
                    .has_background_color = false,
                    .has_background_image = false,
                    .has_can_disable = false,
                    .has_confirm_text = false,
                    .has_draw_arrows = false,
                    .has_disabled = false,
                    .has_hide_action = false,
                    .has_hide_cancel = false,
                    .has_hide_confirm = false,
                    .has_is_header = false,
                    .has_unselectable = false,
                    .has_alignment = false,
                };
                strncpy(state->items[item_index].features.alignment, "left", sizeof(state->items[item_index].features.alignment) - 1);
                strncpy(state->items[item_index].features.confirm_text, confirm_text, sizeof(state->items[item_index].features.confirm_text) - 1);
                if (default_background_image != NULL)
                {
                    strncpy(state->items[item_index].features.background_image, default_background_image, sizeof(state->items[item_index].features.background_image) - 1);
                    if (access(default_background_image, F_OK) != -1)
                    {
                        state->items[item_index].features.background_image_exists = true;
                    }
                }
                if (default_background_color != NULL)
                {
                    strncpy(state->items[item_index].features.background_color, default_background_color, sizeof(state->items[item_index].features.background_color) - 1);
                }

                item_index++;
            }

            if (*line_end == '\0')
            {
                break;
            }
            line_start = line_end + 1;
        }

        free(contents);
        return state;
    }

    JSON_Value *root_value;
    if (strcmp(filename, "-") == 0)
    {
        char *contents = read_stdin();
        if (contents == NULL)
        {
            log_error("Failed to read stdin");
            free(state);
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
        free(state);
        return NULL;
    }

    JSON_Array *items_array;
    if (strlen(item_key) == 0)
    {
        items_array = json_value_get_array(root_value);
    }
    else
    {
        items_array = json_value_get_array(json_object_get_value(json_value_get_object(root_value), item_key));
    }

    size_t item_count = json_array_get_count(items_array);

    state->items = malloc(sizeof(struct ListItem) * item_count);
    state->has_options = false;

    if (strlen(item_key) == 0)
    {
        for (size_t i = 0; i < item_count; i++)
        {
            const char *name = json_array_get_string(items_array, i);
            state->items[i].name = name ? strdup(name) : "";

            // set defaults for the other fields
            state->items[i].has_features = false;
            state->items[i].has_options = false;
            state->items[i].has_selected = false;
            state->items[i].option_count = 0;
            state->items[i].options = NULL;
            state->items[i].selected = 0;
            state->items[i].initial_selected = 0;
            state->items[i].features = (struct ListItemFeature){
                .background_color = "",
                .background_image = "",
                .background_image_exists = false,
                .can_disable = false,
                .confirm_text = "",
                .disabled = false,
                .draw_arrows = false,
                .hide_action = false,
                .hide_cancel = false,
                .hide_confirm = false,
                .is_header = false,
                .unselectable = false,
                .alignment = "",
                .has_background_color = false,
                .has_background_image = false,
                .has_can_disable = false,
                .has_confirm_text = false,
                .has_disabled = false,
                .has_draw_arrows = false,
                .has_hide_action = false,
                .has_hide_cancel = false,
                .has_hide_confirm = false,
                .has_is_header = false,
                .has_unselectable = false,
                .has_alignment = false,
            };
            strncpy(state->items[i].features.alignment, "left", sizeof(state->items[i].features.alignment) - 1);
            strncpy(state->items[i].features.confirm_text, confirm_text, sizeof(state->items[i].features.confirm_text) - 1);
            if (default_background_image != NULL)
            {
                strncpy(state->items[i].features.background_image, default_background_image, sizeof(state->items[i].features.background_image) - 1);
                if (access(default_background_image, F_OK) != -1)
                {
                    state->items[i].features.background_image_exists = true;
                }
            }
            if (default_background_color != NULL)
            {
                strncpy(state->items[i].features.background_color, default_background_color, sizeof(state->items[i].features.background_color) - 1);
            }
        }
    }
    else
    {
        for (size_t i = 0; i < item_count; i++)
        {
            JSON_Object *item = json_array_get_object(items_array, i);

            const char *name = json_object_get_string(item, "name");
            state->items[i].name = name ? strdup(name) : "";

            // read in the options from the json object
            // if there are no options, set the options to an empty array
            // if there are options, treat them as a list of strings
            JSON_Array *options_array = json_object_get_array(item, "options");
            size_t options_count = json_array_get_count(options_array);
            state->items[i].options = malloc(sizeof(char *) * options_count);
            state->items[i].option_count = options_count;
            for (size_t j = 0; j < options_count; j++)
            {
                const char *option = json_array_get_string(options_array, j);
                state->items[i].options[j] = option ? strdup(option) : "";
            }

            if (options_count > 0)
            {
                state->has_options = true;
                state->items[i].has_options = true;
            }
            else
            {
                state->items[i].has_options = false;
            }

            // read in the current option index from the json object
            // if there is no current option index, set it to 0
            // if there is a current option index, treat it as an integer
            if (json_object_has_value(item, "selected"))
            {
                state->items[i].selected = json_object_get_number(item, "selected");
                if (state->items[i].selected < 0)
                {
                    char error_message[256];
                    snprintf(error_message, sizeof(error_message), "Item %s has a selected option index of %d, which is less than 0. Setting to 0.", state->items[i].name, state->items[i].selected);
                    log_error(error_message);
                    state->items[i].selected = 0;
                }
                if (state->items[i].selected >= options_count)
                {
                    char error_message[256];
                    snprintf(error_message, sizeof(error_message), "Item %s has a selected option index of %d, which is greater than the number of options %d. Setting to last option.", state->items[i].name, state->items[i].selected, options_count);
                    log_error(error_message);
                    state->items[i].selected = options_count - 1;
                    if (state->items[i].selected < 0)
                    {
                        state->items[i].selected = 0;
                    }
                }
                state->items[i].has_selected = true;
            }
            else
            {
                state->items[i].selected = 0;
                state->items[i].has_selected = false;
            }

            state->items[i].initial_selected = state->items[i].selected;

            state->items[i].features = (struct ListItemFeature){
                .background_color = "",
                .background_image = "",
                .background_image_exists = false,
                .can_disable = false,
                .confirm_text = "",
                .disabled = false,
                .draw_arrows = false,
                .hide_action = false,
                .hide_cancel = false,
                .hide_confirm = false,
                .is_header = false,
                .unselectable = false,
                .alignment = "",
                .has_background_color = false,
                .has_background_image = false,
                .has_can_disable = false,
                .has_confirm_text = false,
                .has_disabled = false,
                .has_draw_arrows = false,
                .has_hide_action = false,
                .has_hide_cancel = false,
                .has_hide_confirm = false,
                .has_is_header = false,
                .has_unselectable = false,
                .has_alignment = false,
            };
            strncpy(state->items[i].features.alignment, "left", sizeof(state->items[i].features.alignment) - 1);
            strncpy(state->items[i].features.confirm_text, confirm_text, sizeof(state->items[i].features.confirm_text) - 1);
            state->items[i].has_features = false;
            if (json_object_has_value(item, "features"))
            {
                state->items[i].has_features = true;
                JSON_Object *features = json_object_get_object(item, "features");

                // read in the background_image from the json object
                // if there is no background_image, set it to ""
                // if there is a background_image, treat it as a string
                const char *background_image = json_object_get_string(features, "background_image");
                if (background_image != NULL)
                {
                    strncpy(state->items[i].features.background_image, background_image, sizeof(state->items[i].features.background_image) - 1);
                    if (access(background_image, F_OK) != -1)
                    {
                        state->items[i].features.background_image_exists = true;
                    }
                    state->items[i].features.has_background_image = true;
                }
                else
                {
                    if (default_background_image != NULL)
                    {
                        strncpy(state->items[i].features.background_image, default_background_image, sizeof(state->items[i].features.background_image) - 1);
                        if (access(default_background_image, F_OK) != -1)
                        {
                            state->items[i].features.background_image_exists = true;
                        }
                        state->items[i].features.has_background_image = true;
                    }
                    else
                    {
                        state->items[i].features.has_background_image = false;
                    }
                }

                // read in the background_color from the json object
                // if there is no background_color, set it to ""
                // if there is a background_color, treat it as a string
                const char *background_color = json_object_get_string(features, "background_color");
                if (background_color != NULL)
                {
                    strncpy(state->items[i].features.background_color, background_color, sizeof(state->items[i].features.background_color) - 1);
                    state->items[i].features.has_background_color = true;
                }
                else
                {
                    if (default_background_color != NULL)
                    {
                        strncpy(state->items[i].features.background_color, default_background_color, sizeof(state->items[i].features.background_color) - 1);
                        state->items[i].features.has_background_color = true;
                    }
                    else
                    {
                        state->items[i].features.has_background_color = false;
                    }
                }

                // read in the can_disable from the json object
                // if there is no can_disable, set it to false
                // if there is a can_disable, treat it as a boolean
                if (json_object_get_boolean(features, "can_disable") == 1)
                {
                    state->items[i].features.can_disable = true;
                    state->items[i].features.has_can_disable = true;
                }
                else if (json_object_get_boolean(features, "can_disable") == 0)
                {
                    state->items[i].features.can_disable = false;
                    state->items[i].features.has_can_disable = true;
                }
                else
                {
                    state->items[i].features.can_disable = false;
                    state->items[i].features.has_can_disable = false;
                }

                // read in the disabled from the json object
                // if there is no disabled, set it to false
                // if there is an disabled, treat it as a boolean
                if (json_object_get_boolean(features, "disabled") == 1)
                {
                    state->items[i].features.disabled = true;
                    state->items[i].features.has_disabled = true;
                }
                else if (json_object_get_boolean(features, "disabled") == 0)
                {
                    state->items[i].features.disabled = false;
                    state->items[i].features.has_disabled = true;
                    if (!state->items[i].features.can_disable)
                    {
                        char error_message[256];
                        snprintf(error_message, sizeof(error_message), "Item %s has no can_disable, but is disabled", state->items[i].name);
                        log_error(error_message);
                    }
                }
                else
                {
                    state->items[i].features.disabled = false;
                    state->items[i].features.has_disabled = false;
                }

                // read in the draw_arrows from the json object
                // if there is no draw_arrows, set it to false
                // if there is a draw_arrows, treat it as a boolean
                if (json_object_get_boolean(features, "draw_arrows") == 1)
                {
                    state->items[i].features.draw_arrows = true;
                    state->items[i].features.has_draw_arrows = true;
                }
                else if (json_object_get_boolean(features, "draw_arrows") == 0)
                {
                    state->items[i].features.draw_arrows = false;
                    state->items[i].features.has_draw_arrows = true;
                }
                else
                {
                    state->items[i].features.draw_arrows = false;
                    state->items[i].features.has_draw_arrows = false;
                }

                // read in the hide_action from the json object
                // if there is no hide_action, set it to false
                // if there is a hide_action, treat it as a boolean
                if (json_object_get_boolean(features, "hide_action") == 1)
                {
                    state->items[i].features.hide_action = true;
                    state->items[i].features.has_hide_action = true;
                }
                else if (json_object_get_boolean(features, "hide_action") == 0)
                {
                    state->items[i].features.hide_action = false;
                    state->items[i].features.has_hide_action = true;
                }
                else
                {
                    state->items[i].features.hide_action = false;
                    state->items[i].features.has_hide_action = false;
                }

                // read in the hide_cancel from the json object
                // if there is no hide_cancel, set it to false
                // if there is a hide_cancel, treat it as a boolean
                if (json_object_get_boolean(features, "hide_cancel") == 1)
                {
                    state->items[i].features.hide_cancel = true;
                    state->items[i].features.has_hide_cancel = true;
                }
                else if (json_object_get_boolean(features, "hide_cancel") == 0)
                {
                    state->items[i].features.hide_cancel = false;
                    state->items[i].features.has_hide_cancel = true;
                }
                else
                {
                    state->items[i].features.hide_cancel = false;
                    state->items[i].features.has_hide_cancel = false;
                }

                // read in the hide_confirm from the json object
                // if there is no hide_confirm, set it to false
                // if there is a hide_confirm, treat it as a boolean
                if (json_object_get_boolean(features, "hide_confirm") == 1)
                {
                    state->items[i].features.hide_confirm = true;
                    state->items[i].features.has_hide_confirm = true;
                }
                else if (json_object_get_boolean(features, "hide_confirm") == 0)
                {
                    state->items[i].features.hide_confirm = false;
                    state->items[i].features.has_hide_confirm = true;
                }
                else
                {
                    state->items[i].features.hide_confirm = false;
                    state->items[i].features.has_hide_confirm = false;
                }

                // read in the unselectable from the json object
                // if there is no unselectable, set it to false
                // if there is a unselectable, treat it as a boolean
                if (json_object_get_boolean(features, "unselectable") == 1)
                {
                    state->items[i].features.unselectable = true;
                    state->items[i].features.has_unselectable = true;
                }
                else if (json_object_get_boolean(features, "unselectable") == 0)
                {
                    state->items[i].features.unselectable = false;
                    state->items[i].features.has_unselectable = true;
                }
                else
                {
                    state->items[i].features.unselectable = false;
                    state->items[i].features.has_unselectable = false;
                }

                // read in the is_header from the json object
                // if there is no is_header, set it to false
                // if there is a is_header, treat it as a boolean
                // headers are not selectable, so this has to go last such that we can set the unselectable flag
                if (json_object_get_boolean(features, "is_header") == 1)
                {
                    state->items[i].features.is_header = true;
                    state->items[i].features.has_is_header = true;
                    state->items[i].features.unselectable = true;
                }
                else if (json_object_get_boolean(features, "is_header") == 0)
                {
                    state->items[i].features.is_header = false;
                    state->items[i].features.has_is_header = true;
                }
                else
                {
                    state->items[i].features.is_header = false;
                    state->items[i].features.has_is_header = false;
                }

                // read in the alignment from the json object
                // if there is no alignment, set it to 'left'
                // if there is a alignment, it should be 'left', 'center', or 'right'
                const char *alignment = json_object_get_string(features, "alignment");
                if (alignment != NULL)
                {
                    if (strcmp(alignment, "left") == 0 || strcmp(alignment, "center") == 0 || strcmp(alignment, "right") == 0)
                    {
                        strncpy(state->items[i].features.alignment, alignment, sizeof(state->items[i].features.alignment) - 1);
                        state->items[i].features.has_alignment = true;
                    }
                    else
                    {
                        char error_message[256];
                        snprintf(error_message, sizeof(error_message), "Item %s has invalid alignment %s. Must be 'left', 'center', or 'right'. Using default (left).", state->items[i].name, alignment);
                        log_error(error_message);
                        strncpy(state->items[i].features.alignment, "left", sizeof(state->items[i].features.alignment) - 1);
                        state->items[i].features.has_alignment = false;
                    }
                }
                else
                {
                    strncpy(state->items[i].features.alignment, "left", sizeof(state->items[i].features.alignment) - 1);
                    state->items[i].features.has_alignment = false;
                }

                // read in the alignment from the json object
                // if there is no alignment, set it to 'left'
                // if there is a alignment, it should be 'left', 'center', or 'right'
                const char *confirm_text = json_object_get_string(features, "confirm_text");
                if (confirm_text != NULL)
                {
                    if (strlen(confirm_text) > 0)
                    {
                        strncpy(state->items[i].features.confirm_text, confirm_text, sizeof(state->items[i].features.confirm_text) - 1);
                        state->items[i].features.has_confirm_text = true;
                    }
                }
            }
            else
            {
                if (default_background_image != NULL)
                {
                    strncpy(state->items[i].features.background_image, default_background_image, sizeof(state->items[i].features.background_image) - 1);
                    if (access(default_background_image, F_OK) != -1)
                    {
                        state->items[i].features.background_image_exists = true;
                    }
                    state->items[i].features.has_background_image = true;
                }
                else
                {
                    state->items[i].features.has_background_image = false;
                }

                if (default_background_color != NULL)
                {
                    strncpy(state->items[i].features.background_color, default_background_color, sizeof(state->items[i].features.background_color) - 1);
                    state->items[i].features.has_background_color = true;
                }
                else
                {
                    state->items[i].features.has_background_color = false;
                }
            }
        }
    }

    state->last_visible = (item_count < max_row_count) ? item_count : max_row_count;
    state->first_visible = 0;
    state->selected = 0;
    state->item_count = item_count;

    json_value_free(root_value);
    return state;
}

// handle_input interprets input events and mutates app state
void handle_input(struct AppState *state)
{
    // do not redraw by default
    state->redraw = 0;

    if (!state->list_state->items[state->list_state->selected].features.background_image_exists && state->list_state->items[state->list_state->selected].features.background_image != NULL)
    {
        if (access(state->list_state->items[state->list_state->selected].features.background_image, F_OK) != -1)
        {
            state->list_state->items[state->list_state->selected].features.background_image_exists = true;
            state->redraw = 1;
        }
    }

    PAD_poll();

    // discount the title from the max row count
    int max_row_count = MAIN_ROW_COUNT;
    if (strlen(state->title) > 0)
    {
        max_row_count -= 1;
    }

    bool is_action_button_pressed = false;
    bool is_cancel_button_pressed = false;
    bool is_confirm_button_pressed = false;
    bool is_enable_button_pressed = false;
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
        else if (strcmp(state->enable_button, "A") == 0)
        {
            is_enable_button_pressed = true;
        }
    }
    else if (PAD_justReleased(BTN_B))
    {
        if (strcmp(state->action_button, "B") == 0)
        {
            is_action_button_pressed = true;
        }
        else if (strcmp(state->cancel_button, "B") == 0)
        {
            is_cancel_button_pressed = true;
        }
        else if (strcmp(state->confirm_button, "B") == 0)
        {
            is_confirm_button_pressed = true;
        }
        else if (strcmp(state->enable_button, "B") == 0)
        {
            is_enable_button_pressed = true;
        }
    }
    else if (PAD_justReleased(BTN_X))
    {
        if (strcmp(state->action_button, "X") == 0)
        {
            is_action_button_pressed = true;
        }
        else if (strcmp(state->cancel_button, "X") == 0)
        {
            is_cancel_button_pressed = true;
        }
        else if (strcmp(state->confirm_button, "X") == 0)
        {
            is_confirm_button_pressed = true;
        }
        else if (strcmp(state->enable_button, "X") == 0)
        {
            is_enable_button_pressed = true;
        }
    }
    else if (PAD_justReleased(BTN_Y))
    {
        if (strcmp(state->action_button, "Y") == 0)
        {
            is_action_button_pressed = true;
        }
        else if (strcmp(state->cancel_button, "Y") == 0)
        {
            is_cancel_button_pressed = true;
        }
        else if (strcmp(state->confirm_button, "Y") == 0)
        {
            is_confirm_button_pressed = true;
        }
        else if (strcmp(state->enable_button, "Y") == 0)
        {
            is_enable_button_pressed = true;
        }
    }

    if (is_action_button_pressed && !state->list_state->items[state->list_state->selected].features.hide_action)
    {
        state->redraw = 0;
        state->quitting = 1;
        state->exit_code = ExitCodeActionButton;
        return;
    }

    if (is_cancel_button_pressed && !state->list_state->items[state->list_state->selected].features.hide_cancel)
    {
        state->redraw = 0;
        state->quitting = 1;
        state->exit_code = ExitCodeCancelButton;
        return;
    }

    bool force_hide_confirm = false;
    if (state->list_state->items[state->list_state->selected].has_options && state->list_state->items[state->list_state->selected].initial_selected == state->list_state->items[state->list_state->selected].selected)
    {
        force_hide_confirm = true;
    }

    if (state->list_state->items[state->list_state->selected].features.hide_confirm)
    {
        force_hide_confirm = true;
    }

    if (is_confirm_button_pressed && !force_hide_confirm)
    {
        state->redraw = 0;
        state->quitting = 1;
        state->exit_code = ExitCodeSuccess;
        return;
    }

    // if the enable button is pressed, toggle the enabled state of the currently selected item
    if (is_enable_button_pressed)
    {
        if (state->list_state->items[state->list_state->selected].features.can_disable)
        {
            state->redraw = 1;
            state->list_state->items[state->list_state->selected].features.disabled = !state->list_state->items[state->list_state->selected].features.disabled;
        }
        return;
    }

    if (PAD_justReleased(BTN_MENU))
    {
        state->redraw = 0;
        state->quitting = 1;
        state->exit_code = ExitCodeMenuButton;
        return;
    }

    if (PAD_justRepeated(BTN_UP))
    {
        if (state->list_state->selected == 0 && !PAD_justPressed(BTN_UP))
        {
            state->redraw = 0;
        }
        else
        {
            state->list_state->selected -= 1;
            while (state->list_state->items[state->list_state->selected].features.is_header || state->list_state->items[state->list_state->selected].features.unselectable)
            {
                state->list_state->selected -= 1;
                if (state->list_state->selected < 0)
                {
                    break;
                }
            }

            if (state->list_state->selected < 0)
            {
                state->list_state->selected = state->list_state->item_count - 1;
                while (state->list_state->items[state->list_state->selected].features.is_header || state->list_state->items[state->list_state->selected].features.unselectable)
                {
                    state->list_state->selected -= 1;
                }

                int start = state->list_state->item_count - max_row_count;
                state->list_state->first_visible = (start < 0) ? 0 : start;
                state->list_state->last_visible = state->list_state->item_count;
            }
            else if (state->list_state->selected < state->list_state->first_visible)
            {
                state->list_state->first_visible -= 1;
                state->list_state->last_visible -= 1;
            }
            state->redraw = 1;
        }
    }
    else if (PAD_justRepeated(BTN_DOWN))
    {
        if (state->list_state->selected == state->list_state->item_count - 1 && !PAD_justPressed(BTN_DOWN))
        {
            state->redraw = 0;
        }
        else
        {
            state->list_state->selected += 1;
            while (state->list_state->items[state->list_state->selected].features.is_header || state->list_state->items[state->list_state->selected].features.unselectable)
            {
                state->list_state->selected += 1;
                if (state->list_state->selected >= state->list_state->item_count)
                {
                    break;
                }
            }

            if (state->list_state->selected >= state->list_state->item_count)
            {
                state->list_state->selected = 0;
                while (state->list_state->items[state->list_state->selected].features.is_header || state->list_state->items[state->list_state->selected].features.unselectable)
                {
                    state->list_state->selected += 1;
                }

                state->list_state->first_visible = 0;
                state->list_state->last_visible = (state->list_state->item_count < max_row_count) ? state->list_state->item_count : max_row_count;
            }
            else if (state->list_state->selected >= state->list_state->last_visible)
            {
                state->list_state->first_visible += 1;
                state->list_state->last_visible += 1;
            }
            state->redraw = 1;
        }
    }
    else if (PAD_justRepeated(BTN_LEFT))
    {
        // if the state has options, cycle through the options
        if (state->list_state->has_options)
        {
            if (!state->list_state->items[state->list_state->selected].features.disabled)
            {
                state->list_state->items[state->list_state->selected].selected -= 1;
                if (state->list_state->items[state->list_state->selected].selected < 0)
                {
                    state->list_state->items[state->list_state->selected].selected = state->list_state->items[state->list_state->selected].option_count - 1;
                }
            }
        }
        else
        {
            state->list_state->selected -= max_row_count;
            while (state->list_state->items[state->list_state->selected].features.is_header || state->list_state->items[state->list_state->selected].features.unselectable)
            {
                state->list_state->selected -= 1;
            }

            if (state->list_state->selected < 0)
            {
                state->list_state->selected = 0;
                state->list_state->first_visible = 0;
                state->list_state->last_visible = (state->list_state->item_count < max_row_count) ? state->list_state->item_count : max_row_count;
            }
            else if (state->list_state->selected < state->list_state->first_visible)
            {
                state->list_state->first_visible -= max_row_count;
                if (state->list_state->first_visible < 0)
                {
                    state->list_state->first_visible = 0;
                }
                state->list_state->last_visible = state->list_state->first_visible + max_row_count;
            }
        }
        state->redraw = 1;
    }
    else if (PAD_justRepeated(BTN_RIGHT))
    {
        // if the state has options, cycle through the options
        if (state->list_state->has_options)
        {
            if (!state->list_state->items[state->list_state->selected].features.disabled)
            {
                state->list_state->items[state->list_state->selected].selected += 1;
                if (state->list_state->items[state->list_state->selected].selected >= state->list_state->items[state->list_state->selected].option_count)
                {
                    state->list_state->items[state->list_state->selected].selected = 0;
                }
            }
        }
        else
        {
            state->list_state->selected += max_row_count;
            while (state->list_state->items[state->list_state->selected].features.is_header || state->list_state->items[state->list_state->selected].features.unselectable)
            {
                state->list_state->selected += 1;
            }

            if (state->list_state->selected >= state->list_state->item_count)
            {
                state->list_state->selected = state->list_state->item_count - 1;
                int start = state->list_state->item_count - max_row_count;
                state->list_state->first_visible = (start < 0) ? 0 : start;
                state->list_state->last_visible = state->list_state->item_count;
            }
            else if (state->list_state->selected >= state->list_state->last_visible)
            {
                state->list_state->last_visible += max_row_count;
                if (state->list_state->last_visible > state->list_state->item_count)
                {
                    state->list_state->last_visible = state->list_state->item_count;
                }
                state->list_state->first_visible = state->list_state->last_visible - max_row_count;
            }
        }
        state->redraw = 1;
    }
}

// detects if a string is a hex color
bool detect_hex_color(const char *hex)
{
    if (hex[0] != '#')
    {
        return false;
    }

    hex++;
    int r, g, b;
    if (sscanf(hex, "%02x%02x%02x", &r, &g, &b) == 3)
    {
        return true;
    }

    return false;
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

// turns an SDL_Color into a uint32_t
uint32_t sdl_color_to_uint32(SDL_Color color)
{
    return (uint32_t)((color.r << 16) + (color.g << 8) + (color.b << 0));
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

// draw_background draws the background of the list
bool draw_background(SDL_Surface *screen, struct AppState *state)
{
    // render a background color
    char hex_color[1024] = "#000000";
    if (state->list_state->items[state->list_state->selected].features.background_color != NULL)
    {
        strncpy(hex_color, state->list_state->items[state->list_state->selected].features.background_color, sizeof(hex_color));
    }

    SDL_Color background_color = hex_to_sdl_color(hex_color);
    uint32_t color = SDL_MapRGBA(screen->format, background_color.r, background_color.g, background_color.b, 255);
    SDL_FillRect(screen, NULL, color);

    bool should_draw_background_image = false;
    if (state->list_state->items[state->list_state->selected].features.background_image_exists && access(state->list_state->items[state->list_state->selected].features.background_image, F_OK) != -1)
    {
        should_draw_background_image = true;
    }

    // check if there is an image and it is accessible
    if (should_draw_background_image)
    {
        SDL_Surface *surface = IMG_Load(state->list_state->items[state->list_state->selected].features.background_image);
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

    return should_draw_background_image;
}

// draw_screen interprets the app state and draws it to the screen
void draw_screen(SDL_Surface *screen, struct AppState *state, int ow, bool should_draw_background_image)
{
    bool force_hide_confirm = false;
    if (state->list_state->items[state->list_state->selected].has_options && state->list_state->items[state->list_state->selected].initial_selected == state->list_state->items[state->list_state->selected].selected)
    {
        force_hide_confirm = true;
    }

    if (state->list_state->items[state->list_state->selected].features.hide_confirm)
    {
        force_hide_confirm = true;
    }

    // draw the button group on the right
    // only two buttons can be displayed at a time
    if (force_hide_confirm)
    {
        if (!state->list_state->items[state->list_state->selected].features.hide_cancel)
        {
            GFX_blitButtonGroup((char *[]){state->cancel_button, state->cancel_text, NULL}, 1, screen, 1);
        }
    }
    else if (state->list_state->items[state->list_state->selected].features.hide_cancel)
    {
        GFX_blitButtonGroup((char *[]){state->confirm_button, state->list_state->items[state->list_state->selected].features.confirm_text, NULL}, 1, screen, 1);
    }
    else
    {
        GFX_blitButtonGroup((char *[]){state->cancel_button, state->cancel_text, state->confirm_button, state->list_state->items[state->list_state->selected].features.confirm_text, NULL}, 1, screen, 1);
    }

    // if there is a title specified, compute the space needed for it
    int initial_list_y_padding = 0;
    if (strlen(state->title) > 0)
    {
        // Truncate title to avoid battery/wifi icon interference
        int title_available_width = screen->w - SCALE1(PADDING * 3) - ow; // 3 paddings: left, right, and between title and icon pill
        char truncated_title_text[256];
        int title_width = GFX_truncateText(state->fonts.medium, state->title, truncated_title_text, title_available_width, SCALE1(BUTTON_PADDING * 2));

        // compute the x position of the title based on the alignment
        int title_x_pos;
        const char *title_alignment = state->title_alignment ? state->title_alignment : "left";
        if (strcmp(title_alignment, "center") == 0)
        {
            title_x_pos = (screen->w - title_width) / 2 + SCALE1(BUTTON_PADDING);
            int title_interference = title_width - (title_available_width - ow - SCALE1(PADDING)); // extra ow and padding account for centered text, i.e. available width is offset by ow and padding on both sides of screen
            if (title_interference > 0)
            {
                title_x_pos -= title_interference / 2;
            }
        }
        else if (strcmp(title_alignment, "right") == 0)
        {
            title_x_pos = screen->w - title_width - ow - SCALE1(PADDING * 2) + SCALE1(BUTTON_PADDING);
        }
        else // left (default)
        {
            title_x_pos = SCALE1(PADDING + BUTTON_PADDING);
        }

        initial_list_y_padding = PILL_SIZE;
        if (should_draw_background_image)
        {
            int pill_width = MIN(title_available_width, title_width);
            // Calculate pill position based on alignment
            int pill_x_pos;
            if (strcmp(title_alignment, "center") == 0)
            {
                pill_x_pos = (screen->w - pill_width) / 2;
            }
            else if (strcmp(title_alignment, "right") == 0)
            {
                pill_x_pos = screen->w - pill_width - SCALE1(PADDING);
            }
            else // left (default)
            {
                pill_x_pos = SCALE1(PADDING);
            }

            GFX_blitPill(ASSET_BLACK_PILL, screen, &(SDL_Rect){pill_x_pos, SCALE1(PADDING), pill_width, SCALE1(PILL_SIZE)});

            initial_list_y_padding = PILL_SIZE + (PILL_SIZE / 2);
        }

        // draw the title
        SDL_Color text_color = COLOR_GRAY;
        if (should_draw_background_image)
        {
            text_color = COLOR_WHITE;
        }
        SDL_Surface *text = TTF_RenderUTF8_Blended(state->fonts.medium, truncated_title_text, text_color);
        SDL_Rect pos = {
            title_x_pos,
            SCALE1(PADDING + 4),
            text->w,
            text->h};
        SDL_BlitSurface(text, NULL, screen, &pos);
        SDL_FreeSurface(text);
    }

    // the rest of the function is just for drawing your app to the screen
    bool current_item_supports_enabling = false;
    bool current_item_is_enabled = false;
    bool current_item_is_header = false;
    int selected_row = state->list_state->selected - state->list_state->first_visible;
    for (int i = state->list_state->first_visible, j = 0; i < state->list_state->last_visible; i++, j++)
    {
        int available_width = (screen->w) - SCALE1(PADDING * 2);
        bool in_top_row_no_title = (j == 0 && strlen(state->title) == 0);
        // Account for the space taken up by ow and it's padding
        if (in_top_row_no_title)
        {
            available_width -= (ow + SCALE1(PADDING));
        }
        // compute the string representation of the current item
        // to include the current option if there are any options
        // the output should be in the format of:
        // item.name: <selected>
        // if there are no options, the output should be:
        // item.name
        char display_text[256];
        char display_selected_text[256];
        char *alignment = state->list_state->items[i].features.alignment;
        bool is_hex_color = false;
        strncpy(display_selected_text, "", sizeof(display_selected_text));
        if (state->list_state->items[i].option_count > 0)
        {
            char *selected = state->list_state->items[i].options[state->list_state->items[i].selected];
            is_hex_color = detect_hex_color(selected);
            if (strcmp(alignment, "left") == 0)
            {
                snprintf(display_text, sizeof(display_text), "%s", state->list_state->items[i].name);
                if (state->list_state->items[i].features.draw_arrows)
                {
                    snprintf(display_selected_text, sizeof(display_selected_text), "‹ %s ›", selected);
                }
                else
                {
                    snprintf(display_selected_text, sizeof(display_selected_text), "%s", selected);
                }
            }
            else
            {
                if (state->list_state->items[i].features.draw_arrows)
                {
                    snprintf(display_text, sizeof(display_text), "%s: ‹ %s ›", state->list_state->items[i].name, selected);
                }
                else
                {
                    snprintf(display_text, sizeof(display_text), "%s: %s", state->list_state->items[i].name, selected);
                }
            }
        }
        else
        {
            snprintf(display_text, sizeof(display_text), "%s", state->list_state->items[i].name);
        }

        SDL_Color text_color = COLOR_WHITE;
        if (state->list_state->items[i].features.disabled)
        {
            text_color = (SDL_Color){TRIAD_DARK_GRAY};
        }
        if (state->list_state->items[i].features.is_header || state->list_state->items[i].features.unselectable)
        {
            text_color = COLOR_LIGHT_TEXT;
        }

        int color_placeholder_height;
        TTF_SizeUTF8(state->fonts.medium, " ", NULL, &color_placeholder_height);
        int color_box_space = 0;
        if (is_hex_color)
        {
            color_box_space += (color_placeholder_height + SCALE1(PADDING));
            available_width -= color_box_space;
        }

        char truncated_display_text[256];
        int text_width = GFX_truncateText(state->fonts.large, display_text, truncated_display_text, available_width, SCALE1(BUTTON_PADDING * 2));
        int pill_width = MIN(available_width, text_width) + color_box_space;

        if (j == selected_row)
        {
            text_color = COLOR_BLACK;
            current_item_is_enabled = state->list_state->items[i].features.disabled;
            if (state->list_state->items[i].features.disabled)
            {
                text_color = (SDL_Color){TRIAD_LIGHT_GRAY};
            }
            if (state->list_state->items[i].features.is_header || state->list_state->items[i].features.unselectable)
            {
                current_item_is_header = true;
                text_color = COLOR_LIGHT_TEXT;
            }
            if (state->list_state->items[i].features.can_disable)
            {
                current_item_supports_enabling = true;
            }

            // Calculate pill position based on alignment
            int pill_x_pos;
            if (strcmp(alignment, "center") == 0)
            {
                pill_x_pos = (screen->w - pill_width) / 2;
            }
            else if (strcmp(alignment, "right") == 0)
            {
                pill_x_pos = screen->w - pill_width - SCALE1(PADDING);
            }
            else // left (default)
            {
                pill_x_pos = SCALE1(PADDING);
            }

            // Adjust for the pill position in the top row without title
            if (in_top_row_no_title)
            {
                if (strcmp(alignment, "center") == 0)
                {
                    int interference = pill_width - (available_width - ow - SCALE1(PADDING)); // extra ow and padding account for centered text, i.e. available width is offset by ow and padding on both sides of screen
                    if (interference > 0)
                    {
                        pill_x_pos -= interference / 2;
                    }
                }
                else if (strcmp(alignment, "right") == 0)
                {
                    pill_x_pos -= (ow + SCALE1(PADDING));
                }
            }

            if (strcmp(display_selected_text, "") != 0)
            {
                GFX_blitPill(ASSET_DARK_GRAY_PILL, screen, &(SDL_Rect){pill_x_pos, SCALE1(PADDING + (j * PILL_SIZE) + initial_list_y_padding), screen->w - SCALE1(PADDING + BUTTON_MARGIN), SCALE1(PILL_SIZE)});
            }

            GFX_blitPill(ASSET_WHITE_PILL, screen, &(SDL_Rect){pill_x_pos, SCALE1(PADDING + (j * PILL_SIZE) + initial_list_y_padding), pill_width, SCALE1(PILL_SIZE)});
        }

        SDL_Surface *text;
        text = TTF_RenderUTF8_Blended(state->fonts.large, truncated_display_text, text_color);

        // Calculate text position based on alignment
        int text_x_pos;
        int shadow_x_pos;
        if (strcmp(alignment, "center") == 0)
        {
            text_x_pos = (screen->w - text->w - color_box_space) / 2;
            shadow_x_pos = text_x_pos - 2;
        }
        else if (strcmp(alignment, "right") == 0)
        {
            text_x_pos = screen->w - text->w - SCALE1(PADDING + BUTTON_PADDING) - color_box_space;
            shadow_x_pos = screen->w - text->w - SCALE1(2 + PADDING + BUTTON_PADDING) - color_box_space;
        }
        else // left (default)
        {
            text_x_pos = SCALE1(PADDING + BUTTON_PADDING);
            shadow_x_pos = SCALE1(2 + PADDING + BUTTON_PADDING);
        }

        // Adjust for the pill position in the top row without title
        if (in_top_row_no_title)
        {
            if (strcmp(alignment, "center") == 0)
            {
                int interference = pill_width - (available_width - ow - SCALE1(PADDING)); // extra ow and padding account for centered text, i.e. available width is offset by ow and padding on both sides of screen
                if (interference > 0)
                {
                    text_x_pos -= interference / 2;
                }
            }
            else if (strcmp(alignment, "right") == 0)
            {
                text_x_pos -= (ow + SCALE1(PADDING));
            }
        }

        SDL_Rect pos = {
            text_x_pos,
            SCALE1(PADDING + ((i - state->list_state->first_visible) * PILL_SIZE) + initial_list_y_padding + 4),
            text->w,
            text->h};

        // draw the text as a black shadow
        if (should_draw_background_image && j != selected_row)
        {
            // COLOR_BLACK
            SDL_Surface *accent_text;
            accent_text = TTF_RenderUTF8_Blended(state->fonts.large, truncated_display_text, COLOR_BLACK);
            SDL_Rect accent_pos = {
                shadow_x_pos,
                SCALE1(PADDING + ((i - state->list_state->first_visible) * PILL_SIZE) + initial_list_y_padding + 4 + 2),
                accent_text->w,
                accent_text->h};
            SDL_BlitSurface(accent_text, NULL, screen, &accent_pos);
            SDL_FreeSurface(accent_text);
        }

        SDL_BlitSurface(text, NULL, screen, &pos);
        SDL_FreeSurface(text);

        int initial_cube_x_pos = text_x_pos + text->w;

        // draw the selected option text
        if (strcmp(display_selected_text, "") != 0)
        {
            initial_cube_x_pos = screen->w - SCALE1(PADDING + BUTTON_PADDING) - color_box_space;
            if (j != 0 || strlen(state->title) > 0)
            {
                SDL_Color selected_text_color = COLOR_WHITE;
                if (state->list_state->items[i].features.disabled || state->list_state->items[i].features.unselectable)
                {
                    selected_text_color = COLOR_LIGHT_TEXT;
                }
                SDL_Surface *selected_text;
                selected_text = TTF_RenderUTF8_Blended(state->fonts.large, display_selected_text, selected_text_color);
                pos = (SDL_Rect){screen->w - selected_text->w - SCALE1(PADDING + BUTTON_PADDING) - color_box_space, pos.y, selected_text->w, selected_text->h};
                SDL_BlitSurface(selected_text, NULL, screen, &pos);
                SDL_FreeSurface(selected_text);
            }
        }

        if (is_hex_color)
        {
            // get the hex color from the options array
            char *hex_color = state->list_state->items[i].options[state->list_state->items[i].selected];
            SDL_Color current_color = hex_to_sdl_color(hex_color);
            uint32_t color = SDL_MapRGBA(screen->format, current_color.r, current_color.g, current_color.b, 255);

            // Draw outline cube
            uint32_t outline_color = sdl_color_to_uint32(text_color);
            SDL_Rect outline_rect = {
                initial_cube_x_pos + SCALE1(PADDING),
                SCALE1(PADDING + ((i - state->list_state->first_visible) * PILL_SIZE) + initial_list_y_padding + 5), color_placeholder_height,
                color_placeholder_height};
            SDL_FillRect(screen, &(SDL_Rect){outline_rect.x, outline_rect.y, outline_rect.w, outline_rect.h}, outline_color);

            // Draw color cube
            SDL_Rect color_rect = {
                initial_cube_x_pos + SCALE1(PADDING) + 2,
                SCALE1(PADDING + ((i - state->list_state->first_visible) * PILL_SIZE) + initial_list_y_padding + 5) + 2, color_placeholder_height - 4,
                color_placeholder_height - 4};
            SDL_FillRect(screen, &(SDL_Rect){color_rect.x, color_rect.y, color_rect.w, color_rect.h}, color);
        }
    }

    char enable_button_text[256] = "Enable";
    if (current_item_is_enabled)
    {
        strncpy(enable_button_text, "Disable", sizeof(enable_button_text) - 1);
    }

    // draw the button group on the left
    // this should only display the enable button if the current item supports enabling
    // and should only display the action button if it is assigned to a button
    if (current_item_supports_enabling && strcmp(state->enable_button, "") != 0)
    {
        if (strcmp(state->action_button, "") != 0 && !state->list_state->items[state->list_state->selected].features.hide_action)
        {
            GFX_blitButtonGroup((char *[]){state->enable_button, enable_button_text, state->action_button, state->action_text, NULL}, 0, screen, 0);
        }
        else
        {
            GFX_blitButtonGroup((char *[]){state->enable_button, enable_button_text, NULL}, 0, screen, 0);
        }
    }
    else if (strcmp(state->action_button, "") != 0 && !state->list_state->items[state->list_state->selected].features.hide_action)
    {
        GFX_blitButtonGroup((char *[]){state->action_button, state->action_text, NULL}, 0, screen, 0);
    }

    // don't forget to reset the should_redraw flag
    state->redraw = 0;
}

bool open_fonts(struct AppState *state)
{
    if (state->fonts.default_font != NULL)
    {
        // check if the font path is valid
        if (access(state->fonts.default_font, F_OK) == -1)
        {
            log_error("Invalid font path provided");
            return false;
        }
    }

    if (state->fonts.large_font != NULL)
    {
        // check if the font path is valid
        if (access(state->fonts.large_font, F_OK) == -1)
        {
            log_error("Invalid font path provided");
            return false;
        }

        state->fonts.large = TTF_OpenFont(state->fonts.large_font, SCALE1(FONT_LARGE));
        if (state->fonts.large == NULL)
        {
            log_error("Failed to open large font");
            return false;
        }
        TTF_SetFontStyle(state->fonts.large, TTF_STYLE_BOLD);
    }
    else if (state->fonts.default_font != NULL)
    {
        state->fonts.large = TTF_OpenFont(state->fonts.default_font, SCALE1(FONT_LARGE));
        if (state->fonts.large == NULL)
        {
            log_error("Failed to open default font");
            return false;
        }
        TTF_SetFontStyle(state->fonts.large, TTF_STYLE_BOLD);
    }
    else
    {
        state->fonts.large = font.large;
    }

    if (state->fonts.medium_font != NULL)
    {
        // check if the font path is valid
        if (access(state->fonts.medium_font, F_OK) == -1)
        {
            log_error("Invalid font path provided");
            return false;
        }
        state->fonts.medium = TTF_OpenFont(state->fonts.medium_font, SCALE1(FONT_MEDIUM));
        if (state->fonts.medium == NULL)
        {
            log_error("Failed to open medium font");
            return false;
        }
        TTF_SetFontStyle(state->fonts.medium, TTF_STYLE_BOLD);
    }
    else if (state->fonts.default_font != NULL)
    {
        state->fonts.medium = TTF_OpenFont(state->fonts.default_font, SCALE1(FONT_MEDIUM));
        if (state->fonts.medium == NULL)
        {
            log_error("Failed to open default font");
            return false;
        }
        TTF_SetFontStyle(state->fonts.medium, TTF_STYLE_BOLD);
    }
    else
    {
        state->fonts.medium = font.medium;
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
    else
    {
        exit(ExitCodeError);
    }
}

// parse_arguments parses the arguments using getopt and updates the app state
// supports the following flags:
// - --action-button <button> (default: "")
// - --action-text <text> (default: "ACTION")
// - --background-image <path> (default: empty string)
// - --background-color <hex> (default: empty string)
// - --confirm-button <button> (default: "A")
// - --confirm-text <text> (default: "SELECT")
// - --cancel-button <button> (default: "B")
// - --cancel-text <text> (default: "BACK")
// - --enable-button <button> (default: "Y")
// - --disable-auto-sleep (default: false)
// - --font-default <path> (default: empty string)
// - --font-large <path> (default: empty string)
// - --font-medium <path> (default: empty string)
// - --format <format> (default: "json")
// - --hide-hardware-group (default: false)
// - --title <title> (default: empty string)
// - --title-alignment <alignment> (default: "left")
// - --item-key <key> (default: "items")
// - --write-location <location> (default: "-")
// - --write-value <value> (default: "selected")
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
        {"enable-button", required_argument, 0, 'e'},
        {"file", required_argument, 0, 'f'},
        {"font-default", required_argument, 0, 'l'},
        {"font-large", required_argument, 0, 'L'},
        {"font-medium", required_argument, 0, 'M'},
        {"format", required_argument, 0, 'F'},
        {"hide-hardware-group", no_argument, 0, 'H'},
        {"item-key", required_argument, 0, 'K'},
        {"title", required_argument, 0, 't'},
        {"title-alignment", required_argument, 0, 'T'},
        {"write-location", required_argument, 0, 'w'},
        {"write-value", required_argument, 0, 'W'},
        {"disable-auto-sleep", no_argument, 0, 'U'},
        {0, 0, 0, 0}};

    int opt;
    char *font_path_default = NULL;
    char *font_path_large = NULL;
    char *font_path_medium = NULL;
    while ((opt = getopt_long(argc, argv, "a:A:b:B:c:C:d:D:e:f:F:l:L:M:K:t:T:w:W:UH", long_options, NULL)) != -1)
    {
        switch (opt)
        {
        case 'a':
            strncpy(state->action_button, optarg, sizeof(state->action_button) - 1);
            break;
        case 'A':
            strncpy(state->action_text, optarg, sizeof(state->action_text) - 1);
            break;
        case 'b':
            strncpy(state->background_image, optarg, sizeof(state->background_image));
            break;
        case 'B':
            strncpy(state->background_color, optarg, sizeof(state->background_color));
            break;
        case 'c':
            strncpy(state->confirm_button, optarg, sizeof(state->confirm_button) - 1);
            break;
        case 'C':
            strncpy(state->confirm_text, optarg, sizeof(state->confirm_text) - 1);
            break;
        case 'd':
            strncpy(state->cancel_button, optarg, sizeof(state->cancel_button) - 1);
            break;
        case 'D':
            strncpy(state->cancel_text, optarg, sizeof(state->cancel_text) - 1);
            break;
        case 'e':
            strncpy(state->enable_button, optarg, sizeof(state->enable_button) - 1);
            break;
        case 'f':
            strncpy(state->file, optarg, sizeof(state->file) - 1);
            break;
        case 'F':
            strncpy(state->format, optarg, sizeof(state->format) - 1);
            break;
        case 'H':
            state->show_hardware_group = 0;
            break;
        case 'l':
            strncpy(state->fonts.default_font, optarg, sizeof(state->fonts.default_font) - 1);
            break;
        case 'L':
            strncpy(state->fonts.large_font, optarg, sizeof(state->fonts.large_font) - 1);
            break;
        case 'M':
            strncpy(state->fonts.medium_font, optarg, sizeof(state->fonts.medium_font) - 1);
            break;
        case 'K':
            strncpy(state->item_key, optarg, sizeof(state->item_key) - 1);
            break;
        case 't':
            strncpy(state->title, optarg, sizeof(state->title) - 1);
            break;
        case 'T':
            strncpy(state->title_alignment, optarg, sizeof(state->title_alignment) - 1);
            break;
        case 'w':
            strncpy(state->write_location, optarg, sizeof(state->write_location) - 1);
            break;
        case 'W':
            strncpy(state->write_value, optarg, sizeof(state->write_value) - 1);
            break;
        case 'U':
            state->disable_auto_sleep = true;
            break;
        default:
            return false;
        }
    }

    // validate title alignment
    if (strcmp(state->title_alignment, "left") != 0 && strcmp(state->title_alignment, "center") != 0 && strcmp(state->title_alignment, "right") != 0)
    {
        log_error("Invalid title alignment provided. Please provide a value of 'left', 'center', or 'right'.");
        return false;
    }

    if (strcmp(state->format, "") == 0)
    {
        strncpy(state->format, "json", sizeof(state->format) - 1);
    }

    if (strcmp(state->write_value, "") == 0)
    {
        strncpy(state->write_value, "selected", sizeof(state->write_value) - 1);
    }

    // Apply default values for certain buttons and texts
    if (strcmp(state->action_button, "") == 0)
    {
        strncpy(state->action_button, "", sizeof(state->action_button) - 1);
    }

    if (strcmp(state->action_text, "") == 0)
    {
        strncpy(state->action_text, "ACTION", sizeof(state->action_text) - 1);
    }

    if (strcmp(state->cancel_button, "") == 0)
    {
        strncpy(state->cancel_button, "B", sizeof(state->cancel_button) - 1);
    }

    if (strcmp(state->confirm_text, "") == 0)
    {
        strncpy(state->confirm_text, "SELECT", sizeof(state->confirm_text) - 1);
    }

    if (strcmp(state->cancel_text, "") == 0)
    {
        strncpy(state->cancel_text, "BACK", sizeof(state->cancel_text) - 1);
    }

    if (strcmp(state->enable_button, "") == 0)
    {
        strncpy(state->enable_button, "Y", sizeof(state->enable_button) - 1);
    }

    // validate that hardware buttons aren't assigned to more than once
    bool a_button_assigned = false;
    bool b_button_assigned = false;
    bool x_button_assigned = false;
    bool y_button_assigned = false;
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
    if (strcmp(state->enable_button, "A") == 0)
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
    if (strcmp(state->enable_button, "B") == 0)
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
    if (strcmp(state->enable_button, "X") == 0)
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
    if (strcmp(state->enable_button, "Y") == 0)
    {
        if (y_button_assigned)
        {
            log_error("Y button cannot be assigned to more than one button");
            return false;
        }
        y_button_assigned = true;
    }

    // validate that the confirm and cancel buttons are valid
    if (strcmp(state->confirm_button, "A") != 0 && strcmp(state->confirm_button, "B") != 0 && strcmp(state->confirm_button, "X") != 0 && strcmp(state->confirm_button, "Y") != 0)
    {
        log_error("Invalid confirm button provided");
        return false;
    }
    if (strcmp(state->cancel_button, "A") != 0 && strcmp(state->cancel_button, "B") != 0 && strcmp(state->cancel_button, "X") != 0 && strcmp(state->cancel_button, "Y") != 0)
    {
        log_error("Invalid cancel button provided");
        return false;
    }

    if (strlen(state->file) == 0)
    {
        log_error("No input provided");
        return false;
    }

    if (strlen(state->format) == 0)
    {
        log_error("No format provided");
        return false;
    }
    // validate format, and only allow json or text
    if (strcmp(state->format, "json") != 0 && strcmp(state->format, "text") != 0)
    {
        log_error("Invalid format provided");
        return false;
    }

    return true;
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

// write_to_file writes some text to a file
int write_to_file(const char *filename, const char *text)
{
    FILE *file = fopen(filename, "w");
    if (file == NULL)
    {
        log_error("Failed to open write location");
        return ExitCodeError;
    }

    int num_elements = strlen(text) / sizeof(text[0]);
    fwrite(text, sizeof(char), num_elements, file);
    fclose(file);

    return ExitCodeSuccess;
}

// write_output writes the final text to the write location
int write_output(struct AppState *state)
{
    if (strcmp(state->write_value, "selected") == 0)
    {
        if (state->exit_code != ExitCodeSuccess && state->exit_code != ExitCodeActionButton)
        {
            return state->exit_code;
        }

        if (strcmp(state->write_location, "-") != 0)
        {
            log_info(state->list_state->items[state->list_state->selected].name);
        }

        write_to_file(state->write_location, state->list_state->items[state->list_state->selected].name);
        return state->exit_code;
    }

    JSON_Value *root_value = json_value_init_object();
    JSON_Object *root_object = json_value_get_object(root_value);
    char *serialized_string = NULL;
    json_object_set_number(root_object, "selected", state->list_state->selected);

    JSON_Array *items = json_array(json_value_init_array());
    for (int i = 0; i < state->list_state->item_count; i++)
    {
        JSON_Value *val = json_value_init_object();
        JSON_Object *obj = json_value_get_object(val);

        JSON_Value *features_val = json_value_init_object();
        JSON_Object *features = json_value_get_object(features_val);
        if (json_object_dotset_string(obj, "name", state->list_state->items[i].name) == JSONFailure)
        {
            log_error("Failed to set name");
            return ExitCodeSerializeError;
        }

        if (state->list_state->items[i].features.has_alignment)
        {
            if (json_object_dotset_string(features, "alignment", state->list_state->items[i].features.alignment) == JSONFailure)
            {
                log_error("Failed to set alignment");
                return ExitCodeSerializeError;
            }
        }

        if (state->list_state->items[i].features.has_confirm_text)
        {
            if (json_object_dotset_string(features, "confirm_text", state->list_state->items[i].features.confirm_text) == JSONFailure)
            {
                log_error("Failed to set confirm_text");
                return ExitCodeSerializeError;
            }
        }

        if (state->list_state->items[i].features.has_can_disable)
        {
            if (json_object_dotset_boolean(features, "can_disable", state->list_state->items[i].features.can_disable) == JSONFailure)
            {
                log_error("Failed to set can_disable");
                return ExitCodeSerializeError;
            }
        }

        if (state->list_state->items[i].features.has_disabled || state->list_state->items[i].features.has_can_disable)
        {
            if (json_object_dotset_boolean(features, "disabled", state->list_state->items[i].features.disabled))
            {
                log_error("Failed to set enabled");
                return ExitCodeSerializeError;
            }
        }

        if (state->list_state->items[i].features.has_draw_arrows)
        {
            if (json_object_dotset_boolean(features, "draw_arrows", state->list_state->items[i].features.draw_arrows))
            {
                log_error("Failed to set draw_arrows");
                return ExitCodeSerializeError;
            }
        }

        if (state->list_state->items[i].features.has_hide_action)
        {
            if (json_object_dotset_boolean(features, "hide_action", state->list_state->items[i].features.hide_action))
            {
                log_error("Failed to set hide_action");
                return ExitCodeSerializeError;
            }
        }

        if (state->list_state->items[i].features.has_hide_cancel)
        {
            if (json_object_dotset_boolean(features, "hide_cancel", state->list_state->items[i].features.hide_cancel))
            {
                log_error("Failed to set hide_cancel");
                return ExitCodeSerializeError;
            }
        }

        if (state->list_state->items[i].features.has_hide_confirm)
        {
            if (json_object_dotset_boolean(features, "hide_confirm", state->list_state->items[i].features.hide_confirm))
            {
                log_error("Failed to set hide_confirm");
                return ExitCodeSerializeError;
            }
        }

        if (state->list_state->items[i].features.has_is_header)
        {
            if (json_object_dotset_boolean(features, "is_header", state->list_state->items[i].features.is_header))
            {
                log_error("Failed to set is_header");
                return ExitCodeSerializeError;
            }
        }

        if (state->list_state->items[i].features.has_unselectable)
        {
            if (json_object_dotset_boolean(features, "unselectable", state->list_state->items[i].features.unselectable))
            {
                log_error("Failed to set unselectable");
                return ExitCodeSerializeError;
            }
        }

        if (state->list_state->items[i].has_options)
        {
            if (json_object_dotset_number(obj, "selected", state->list_state->items[i].selected) == JSONFailure)
            {
                log_error("Failed to set selected");
                return ExitCodeSerializeError;
            }

            JSON_Array *options = json_array(json_value_init_array());
            for (int j = 0; j < state->list_state->items[i].option_count; j++)
            {
                JSON_Value *option = json_value_init_string(state->list_state->items[i].options[j]);
                if (json_array_append_value(options, option) == JSONFailure)
                {
                    log_error("Failed to append option");
                    return ExitCodeSerializeError;
                }
            }
            if (json_object_dotset_value(obj, "options", json_array_get_wrapping_value(options)) == JSONFailure)
            {
                log_error("Failed to set options");
                return ExitCodeSerializeError;
            }
        }

        // this should always go last
        if (state->list_state->items[i].has_features)
        {
            JSON_Value *features_value = json_object_get_wrapping_value(features);
            if (json_object_dotset_value(obj, "features", features_value) == JSONFailure)
            {
                log_error("Failed to set features");
                return ExitCodeSerializeError;
            }
        }

        JSON_Value *item_value = json_object_get_wrapping_value(obj);
        if (json_array_append_value(items, item_value) == JSONFailure)
        {
            log_error("Failed to append item");
            return ExitCodeSerializeError;
        }
    }

    JSON_Value *items_value = json_array_get_wrapping_value(items);
    if (json_object_dotset_value(root_object, state->item_key, items_value) == JSONFailure)
    {
        log_error("Failed to set items");
        return ExitCodeSerializeError;
    }

    root_value = json_object_get_wrapping_value(root_object);

    serialized_string = json_serialize_to_string_pretty(root_value);
    if (serialized_string == NULL)
    {
        log_error("Failed to serialize");
        return ExitCodeSerializeError;
    }

    if (strcmp(state->write_location, "-") == 0)
    {
        log_info(serialized_string);
    }
    else
    {
        write_to_file(state->write_location, serialized_string);
    }

    json_free_serialized_string(serialized_string);
    json_value_free(root_value);

    return state->exit_code;
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
    char default_enable_button[1024] = "Y";
    char default_confirm_button[1024] = "A";
    char default_confirm_text[1024] = "SELECT";
    char default_file[1024] = "";
    char default_format[1024] = "json";
    char default_item_key[1024] = "";
    char default_write_value[1024] = "selected";
    char default_title[1024] = "";
    char default_title_alignment[1024] = "left";
    char default_write_location[1024] = "-";
    struct AppState state = {
        .exit_code = ExitCodeSuccess,
        .quitting = 0,
        .redraw = 1,
        .show_hardware_group = 1,
        .show_brightness_setting = 0,
        .disable_auto_sleep = false,
        .fonts = {
            .large = NULL,
            .medium = NULL,
            .default_font = NULL,
            .large_font = NULL,
            .medium_font = NULL,
        },
        .list_state = NULL};

    // assign the default values to the app state
    strncpy(state.action_button, default_action_button, sizeof(state.action_button) - 1);
    strncpy(state.action_text, default_action_text, sizeof(state.action_text) - 1);
    strncpy(state.background_image, default_background_image, sizeof(state.background_image));
    strncpy(state.background_color, default_background_color, sizeof(state.background_color));
    strncpy(state.cancel_button, default_cancel_button, sizeof(state.cancel_button) - 1);
    strncpy(state.cancel_text, default_cancel_text, sizeof(state.cancel_text) - 1);
    strncpy(state.confirm_button, default_confirm_button, sizeof(state.confirm_button) - 1);
    strncpy(state.confirm_text, default_confirm_text, sizeof(state.confirm_text) - 1);
    strncpy(state.enable_button, default_enable_button, sizeof(state.enable_button) - 1);
    strncpy(state.file, default_file, sizeof(state.file) - 1);
    strncpy(state.format, default_format, sizeof(state.format) - 1);
    strncpy(state.item_key, default_item_key, sizeof(state.item_key) - 1);
    strncpy(state.write_value, default_write_value, sizeof(state.write_value) - 1);
    strncpy(state.title, default_title, sizeof(state.title) - 1);
    strncpy(state.title_alignment, default_title_alignment, sizeof(state.title_alignment) - 1);
    strncpy(state.write_location, default_write_location, sizeof(state.write_location) - 1);

    // parse the arguments
    if (!parse_arguments(&state, argc, argv))
    {
        return ExitCodeError;
    }

    state.list_state = ListState_New(state.file, state.format, state.item_key, state.title, state.confirm_text, state.background_image, state.background_color);
    if (state.list_state == NULL)
    {
        log_error("Failed to create list state");
        return ExitCodeError;
    }

    if (state.list_state->item_count > 0)
    {
        // if there are items in the list,
        // validate that at least one item is not a header and is selectable
        bool has_selectable = false;
        for (size_t i = 0; i < state.list_state->item_count; i++)
        {
            state.list_state->selected = i;
            if (!state.list_state->items[i].features.is_header && !state.list_state->items[i].features.unselectable)
            {
                has_selectable = true;
                break;
            }
        }
        if (!has_selectable)
        {
            log_error("No selectable items found");
            return ExitCodeError;
        }
    }

    // swallow all stdout from init calls
    // MinUI will sometimes randomly log to stdout
    swallow_stdout_from_function(init);

    signal(SIGINT, signal_handler);
    signal(SIGTERM, signal_handler);

    if (!open_fonts(&state))
    {
        log_error("Failed to open fonts");
        return ExitCodeError;
    }

    // get initial wifi state
    int was_online = PLAT_isOnline();

    // draw the screen at least once
    // handle_input sets state.redraw to 0 if no key is pressed
    int was_ever_drawn = 0;

    if (state.disable_auto_sleep)
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
        PWR_update(&state.redraw, NULL, NULL, NULL);
        bool power_redraw = false;
        if (state.redraw)
        {
            power_redraw = true;
        }

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

        // force a redraw if the screen was never drawn
        if (!was_ever_drawn && !state.redraw)
        {
            state.redraw = 1;
            was_ever_drawn = 1;
        }

        // force a redraw if the power state changed
        if (power_redraw)
        {
            state.redraw = 1;
        }

        // redraw the screen if there has been a change
        if (state.redraw)
        {
            // clear the screen at the beginning of each loop
            GFX_clear(screen);

            bool should_draw_background_image = draw_background(screen, &state);

            int ow = 0;
            if (state.show_hardware_group)
            {
                // draw the hardware information in the top-right
                ow = GFX_blitHardwareGroup(screen, state.show_brightness_setting);

                // draw the setting hints
                if (state.show_brightness_setting && !GetHDMI())
                {
                    GFX_blitHardwareHints(screen, state.show_brightness_setting);
                }
                else
                {
                    GFX_blitButtonGroup((char *[]){BTN_SLEEP == BTN_POWER ? "POWER" : "MENU", "SLEEP", NULL}, 0, screen, 0);
                }
            }

            // your draw logic goes here
            draw_screen(screen, &state, ow, should_draw_background_image);

            // Takes the screen buffer and displays it on the screen
            GFX_flip(screen);
        }
        else
        {
            // Slows down the frame rate to match the refresh rate of the screen
            // when the screen is not being redrawn
            GFX_sync();
        }
    }

    int exit_code = write_output(&state);
    if (exit_code != ExitCodeSuccess)
    {
        return exit_code;
    }

    swallow_stdout_from_function(destruct);

    // exit the program
    return state.exit_code;
}
