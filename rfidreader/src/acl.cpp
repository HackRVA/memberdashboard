#include <Arduino.h>
#include "LittleFS.h"
#include "Hash.h"

#include "acl.h"

#define ACL_FILE "/acl"

void acl_init()
{
    bool success = LittleFS.begin();

    if (success)
    {
        Serial.println("File system mounted with success");
    }
    else
    {
        Serial.println("Error mounting the file system");
        return;
    }
}

void write_acl(String acl[], uint8_t arrsize)
{
    File file = LittleFS.open(ACL_FILE, "w+");

    if (!file)
    {
        Serial.println("Error opening file for writing");
        return;
    }

    Serial.printf("checking the zero index: %s\n", acl[0].c_str());

    for (uint8_t n = 0; n < arrsize; n++)
    {
        if (!acl[n])
            break;
        /* check to see if we are adding zero values to the acl */
        if (acl[n].length() == 0)
        {
            /* skip any zero values */
            continue;
        }

        Serial.printf("writing: %s\n", acl[n].c_str());
        file.println(acl[n]);
    }

    file.close();
}

/**
 * find_id
 * looks for an id in the ACL
 * if it finds the id, it will return true
 * if it doesn't, it will return false
 */
bool find_id(String uid)
{
    File file = LittleFS.open(ACL_FILE, "r");

    Serial.printf("\nlooking for: %s\n\n", uid.c_str());

    if (!file)
    {
        Serial.println("Error opening file for reading");
        return false;
    }

    bool found = false;

    found = file.find(uid.c_str(), uid.length());
    if (found)
    {
        file.close();
        return true;
    }

    file.close();

    return false;
}

String acl_hash()
{
    File file = LittleFS.open(ACL_FILE, "r");

    if (!file)
    {
        Serial.println("Error opening file for reading");
    }

    String hash = sha1(file.readString());
    Serial.println(hash);

    file.close();

    return hash;
}
