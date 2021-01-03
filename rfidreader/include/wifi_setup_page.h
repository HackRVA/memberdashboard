
// Creates an admin page on a webserver which allows the user to update the SSID and Password
// Performs basic checks to ensure that the input values are valid

#ifndef AdminPage_h
#define AdminPage_h
#include <ESP8266WebServer.h>

//Holds the admin webpage in the program memory
const char adminPage[] PROGMEM =
    "<html>"
    "  <head>"
    "    <style>"
    "      input {"
    "        font-size: 1.2em;"
    "        width: 100%;"
    "        max-width: 350px;"
    "        display: block;"
    "        margin: 5px auto;"
    "      }"
    "    </style>"
    "    <script>"
    "        const ssid = document.getElementById('ssid')"
    "        const password = document.getElementById('password')"
    "        const submit = document.getElementById('submit')"
    ""
    "        submit.addEventListener('click', () => {"
    "            fetch('/setup-wifi', {"
    "                method: 'POST',"
    "                options: JSON.stringify({"
    "                    ssid: ssid.value,"
    "                    password: password.value"
    "                })"
    "            })"
    "        })"
    "    </script>"
    "  </head>"
    "  <body>"
    "    <form id='form' action='/admin' method='post'>"
    "      <input"
    "        id='ssid'"
    "        name='ssid'"
    "        type='text'"
    "        minlength='3'"
    "        maxlength='16'"
    "        placeholder='Enter SSID'"
    "      />"
    "      <input"
    "        id='password'"
    "        name='newpassword'"
    "        id='password'"
    "        type='password'"
    "        minlength='3'"
    "        maxlength='16'"
    "        placeholder='Enter Password'"
    "      />"
    "      <input id='submit' type='submit' value='Update' />"
    "    </form>"
    "  </body>"
    "</html>";

#endif
