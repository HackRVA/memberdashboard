#include <Arduino.h>
#include <ArduinoJson.h>
#include <ESP8266WiFi.h>
#include <WiFiClient.h>
#include <ESP8266WebServer.h>

#include "acl.h"
#include "config.h"
#include "wifi_setup_page.h"

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

void serveWifiSetup()
{
    String s = adminPage;             //Read HTML contents
    server.send(200, "text/html", s); //Send web page
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
    config_t *conf = new config_t();
    read_config(conf);

    if (strlen(conf->ssid) == 0 || strlen(conf->password) == 0)
    {
        WiFi.mode(WIFI_AP);
        WiFi.softAP("RFIDReaderConfig", "", 1, false, 1);
        server.on("/", HTTP_GET, serveWifiSetup);
        server.on("/setup-wifi", HTTP_POST, handleWifiSetup);

        Serial.print("IP address: ");
        Serial.println(WiFi.localIP());
        // Set not found response
        server.onNotFound(handleNotFound);
        // Start server
        server.begin();
        Serial.println("HTTP server started");
        return;
    }

    WiFi.mode(WIFI_STA);
    WiFi.begin(conf->ssid, conf->password);
    Serial.println("");

    // Wait for connection
    while (WiFi.status() != WL_CONNECTED)
    {
        delay(500);
        Serial.print(".");
    }
    Serial.println("");
    Serial.print("Connected to ");
    Serial.println(conf->ssid);
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

void handleWifiSetup()
{
    if (server.method() != HTTP_POST)
    {
        server.send(405, "text/plain", "Method Not Allowed");
    }
    else
    {
        StaticJsonDocument<256> doc;
        deserializeJson(doc, server.arg("plain"));
        JsonObject obj = doc.as<JsonObject>();

        if (!doc.containsKey("ssid") && !doc.containsKey("password"))
        {
            server.send(500, "application/json", "{\"error\": \"no ssid and/or password in request \"}");
        }
        write_config(obj["ssid"], obj["password"]);
        server.send(200, "application/json", server.arg("plain"));
        setup_wifi();
    }
}

void wifi_loop(void)
{
    server.handleClient();
}
