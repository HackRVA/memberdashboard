#include <Arduino.h>
#include <ArduinoJson.h>
#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>

#include "acl.h"

// Replace with your network credentials
const char *ssid = "";
const char *password = "";

ESP8266WebServer server(80);

void getACLHash()
{
    server.send(200, "application/json", "{\"acl\": \"" + acl_hash() + "\", \"version\": " + "\"" + "0.0.0" + "\"}");
}

/**
 * handleUpdateACL
 * a POST request that writes a new ACL to memory
 * The endpoint expects JSON with a key called "acl" that is an array of rfid tags
 * ``` {"acl": [2755459513, 848615840]} ```
 */
void handleUpdateACL()
{
    if (server.method() != HTTP_POST)
    {
        server.send(405, "text/plain", "Method Not Allowed");
    }
    else
    {
        StaticJsonDocument<256> doc;
        deserializeJson(doc, server.arg("plain"));

        if (doc.containsKey("acl"))
        {
            JsonArray array = doc["acl"].as<JsonArray>();

            if (array.size() > MAXIMUM_ACL_SIZE || array.size() <= 0)
            {
                server.send(500, "application/json", "{\"error\": \"not a valid access list\"}");
            }

            String new_acl[array.size()];
            uint8_t i = 0;
            for (JsonVariant v : array)
            {
                new_acl[i] = v.as<String>();
                i++;
            }

            write_acl(new_acl, array.size());
        }
        server.send(200, "application/json", server.arg("plain"));
    }
}

/**
 * handleClearACL
 * replace the ACL with an empty ACL
 * this would remove access for everyone
 */
void handleClearACL()
{
    String acl[MAXIMUM_ACL_SIZE] = {};
    write_acl(acl, 0);
    server.send(200, "application/json", "{\"acl\": \"" + acl_hash() + "\", \"version\": " + "\"" + "0.0.0" + "\"}");
}

// Define routing
void restServerRouting()
{
    server.on("/", HTTP_GET, getACLHash);
    server.on(F("/update"), HTTP_POST, handleUpdateACL);
    server.on(F("/clear"), HTTP_POST, handleClearACL);
}

// Manage not found URL
void handleNotFound()
{
    String message = "File Not Found\n\n";
    message += "URI: ";
    message += server.uri();
    message += "\nMethod: ";
    message += (server.method() == HTTP_GET) ? "GET" : "POST";
    message += "\nArguments: ";
    message += server.args();
    message += "\n";
    for (uint8_t i = 0; i < server.args(); i++)
    {
        message += " " + server.argName(i) + ": " + server.arg(i) + "\n";
    }
    server.send(404, "text/plain", message);
}

void setup_wifi(void)
{
    WiFi.mode(WIFI_STA);
    WiFi.begin(ssid, password);
    Serial.println("");

    // Wait for connection
    while (WiFi.status() != WL_CONNECTED)
    {
        delay(500);
        Serial.print(".");
    }
    Serial.println("");
    Serial.print("Connected to ");
    Serial.println(ssid);
    Serial.print("IP address: ");
    Serial.println(WiFi.localIP());

    // Set server routing
    restServerRouting();
    // Set not found response
    server.onNotFound(handleNotFound);
    // Start server
    server.begin();
    Serial.println("HTTP server started");
}

void wifi_loop(void)
{
    server.handleClient();
}
