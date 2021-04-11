#include <Arduino.h>
#include <MFRC522.h>

#include "acl.h"
#include "relay.h"

/**
 * granted
 * Trigger the relay to open and do some user 
 * feedback to let them know that they are in
 * 
 * then relock the door
 */
void granted(uint16_t setDelay)
{
    Serial.println(F("Welcome, You shall pass"));
#if 0
    digitalWrite(blueLed, LED_OFF);   // Turn off blue LED
    digitalWrite(redLed, LED_OFF);  // Turn off red LED
    digitalWrite(greenLed, LED_ON);   // Turn on green LED
#endif
    digitalWrite(RELAY_PIN, HIGH);     // Unlock door!
    delay(setDelay); // Hold door lock open for given seconds
    digitalWrite(RELAY_PIN, LOW);    // Relock door
#if 0
#endif
    delay(1000); // Hold green LED on for a second
}

/**
 * denied
 * do some user feedback to let them 
 * know that they do not have access
 */
void denied()
{
    Serial.println(F("You shall not pass"));
}

#define DEBUG 1
/**
 * getID
 * Get PICC's UID
 */
String getID(MFRC522 reader)
{
    String uid = "";
    /* There are Mifare PICCs which have 4 byte or 7 byte UID care if you use 7 byte PICC
     I think we should assume every PICC as they have 4 byte UID
     Until we support 7 byte PICCs */
    for (uint8_t i = 0; i < reader.uid.size; i++)
    {
        uid += String(reader.uid.uidByte[i], HEX);
    }

    return uid;
}

/**
 * checkID
 * waits for a successful read then it determines if the
 * id exists in the access list
 */
void checkID(MFRC522 reader)
{
    String successRead;

    do
    {
        successRead = getID(reader);
    } while ((successRead.length() / 2) != reader.uid.size);

    if (find_id(successRead))
    {                 // look for the ID in flash
        granted(RELAY_GRANT_DELAY); // Open the door lock for 300 ms
    }
    else
    { // If not, show that the ID was not valid
        denied();
    }
}
