#include <Arduino.h>
#include <ArduinoJson.h>
#include "LittleFS.h"

#define CONFIG_FILE "/acl"
#define MAX_STR_LEN 255
#define SSID_DEFAULT ""
#define PASS_DEFAULT ""

struct config_t
{
    char ssid[MAX_STR_LEN];
    char password[MAX_STR_LEN];
};

// write a config file that contains wifi creds
// read wifi creds from config file
// an endpoint where we can set wifi creds
void write_config(String ssid, String password)
{
    File file = LittleFS.open(CONFIG_FILE, "w+");

    if (!file)
    {
        Serial.println("Error opening file for writing");
        return;
    }

    String json = "{\"ssid\": \"" + ssid + "\", \"password\": \"" + password + "\"}";

    file.println(json);

    file.close();
}

void read_config(config_t *conf)
{
    File file = LittleFS.open(CONFIG_FILE, "r");

    if (!file)
    {
        Serial.println("Error opening file for reading");
        return;
    }
    StaticJsonDocument<256> doc;
    deserializeJson(doc, file);

    if (doc.containsKey("ssid"))
    {
        strlcpy(conf->ssid,                 // <- destination
                doc["ssid"] | SSID_DEFAULT, // <- source
                sizeof(conf->ssid));
        strlcpy(conf->password,                 // <- destination
                doc["password"] | PASS_DEFAULT, // <- source
                sizeof(conf->password));

        file.close();
    }
}
