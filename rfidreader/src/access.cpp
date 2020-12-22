#include <Arduino.h>
#include <MFRC522.h>

#include "acl.h"

byte readCard[4];

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
      digitalWrite(relay, LOW);     // Unlock door!
#endif
    delay(setDelay); // Hold door lock open for given seconds
#if 0
      digitalWrite(relay, HIGH);    // Relock door
#endif
    delay(1000);     // Hold green LED on for a second
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

/**
 * getID
 * Get PICC's UID
 */
uint8_t getID(MFRC522 reader)
{
    /* There are Mifare PICCs which have 4 byte or 7 byte UID care if you use 7 byte PICC
     I think we should assume every PICC as they have 4 byte UID
     Until we support 7 byte PICCs */
    for (uint8_t i = 0; i < 4; i++)
    {
        readCard[i] = reader.uid.uidByte[i];
        Serial.print(readCard[i], HEX);
    }
    Serial.println("");
    reader.PICC_HaltA(); // Stop reading
    return 1;
}

/**
 * checkID
 * waits for a successful read then it determines if the
 * id exists in the access list
 */
void checkID(MFRC522 reader)
{
    uint8_t successRead; // Variable integer to keep if we have Successful Read from Reader

    do
    {
        successRead = getID(reader);

    } while (!successRead);

    if (find_id(readCard))
    { // look for the ID in flash
        granted(300); // Open the door lock for 300 ms
    }
    else
    { // If not, show that the ID was not valid
        denied();
    }
}
