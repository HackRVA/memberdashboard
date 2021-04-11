#include <Arduino.h>
#include <SPI.h>
#include <MFRC522.h>

#include "access.h"

#define RST_PIN 5 // Configurable, see typical pin layout above
#define SS_PIN 4  // Configurable, see typical pin layout above

void ShowReaderDetails(MFRC522 reader)
{
    // Get the MFRC522 software version
    byte v = reader.PCD_ReadRegister(reader.VersionReg);
    Serial.print(F("MFRC522 Software Version: 0x"));
    Serial.print(v, HEX);
    if (v == 0x91)
        Serial.print(F(" = v1.0"));
    else if (v == 0x92)
        Serial.print(F(" = v2.0"));
    else
        Serial.print(F(" (unknown),probably a chinese clone?"));
    Serial.println("");
    // When 0x00 or 0xFF is returned, communication probably failed
    if ((v == 0x00) || (v == 0xFF))
    {
        Serial.println(F("WARNING: Communication failure, is the MFRC522 properly connected?"));
        Serial.println(F("SYSTEM HALTED: Check connections."));
        // Visualize system is halted
        // digitalWrite(greenLed, LED_OFF);  // Make sure green LED is off
        // digitalWrite(blueLed, LED_OFF);   // Make sure blue LED is off
        // digitalWrite(redLed, LED_ON);   // Turn on red LED
        while (true); // do not go further
    }
}

MFRC522 setup_rfid_reader()
{
    MFRC522 reader(SS_PIN, RST_PIN);  // Create MFRC522 instance
    SPI.begin();                      // Init SPI bus
    reader.PCD_Init();                // Init MFRC522
    delay(4);                         // Optional delay. Some board do need more time after init to be ready, see Readme
    reader.PCD_DumpVersionToSerial(); // Show details of PCD - MFRC522 Card Reader details
    Serial.println(F("Scan PICC to see UID, SAK, type, and data blocks..."));

    ShowReaderDetails(reader);

    return reader;
}

void read_rfid(MFRC522 reader)
{
    // Reset the loop if no new card present on the sensor/reader. This saves the entire process when idle.
    if (!reader.PICC_IsNewCardPresent())
    {
        return;
    }

    // Select one of the cards
    if (!reader.PICC_ReadCardSerial())
    {
        return;
    }

    /* check that the user has access */
    checkID(reader);
}
