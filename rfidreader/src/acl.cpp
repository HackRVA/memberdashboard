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

void write_acl(unsigned long acl[])
{
    File file = LittleFS.open(ACL_FILE, "w+");

    if (!file)
    {
        Serial.println("Error opening file for writing");
        return;
    }

    for (uint8_t n = 0; n < MAXIMUM_ACL_SIZE; n++)
    {
        /* check to see if we are adding zero values to the acl */
        if (acl[n] == 0)
        {
            /* skip any zero values */
            continue;
        }

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
bool find_id(byte id[])
{
    File file = LittleFS.open(ACL_FILE, "r");

    unsigned long number = *((unsigned long *)id);
    Serial.printf("\nlooking for: %lu\n\n", number);

    char buf[16];
    ultoa(number, buf, 10);

    if (!file)
    {
        Serial.println("Error opening file for reading");
        return false;
    }

    bool found = false;

    found = file.find(buf);
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
