
#include <Arduino.h>

#define MAX_STR_LEN 255

struct config_t
{
    char ssid[MAX_STR_LEN];
    char password[MAX_STR_LEN];
};

extern void read_config(config_t *conf);
extern void write_config(String ssid, String password);
