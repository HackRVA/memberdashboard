#include <Arduino.h>
#include <MFRC522.h>

#define RST_PIN 5 // Configurable, see typical pin layout above
#define SS_PIN 4  // Configurable, see typical pin layout above

MFRC522 mfrc522;

MFRC522 setup_rfid_reader()
{
	MFRC522 reader(SS_PIN, RST_PIN);  // Create MFRC522 instance
	SPI.begin();					  // Init SPI bus
	reader.PCD_Init();				  // Init MFRC522
	delay(4);						  // Optional delay. Some board do need more time after init to be ready, see Readme
	// reader.PCD_DumpVersionToSerial(); // Show details of PCD - MFRC522 Card Reader details
	// Serial.println(F("Scan PICC to see UID, SAK, type, and data blocks..."));


	return reader;
}

void setup()
{
	Serial.begin(9600); // Initialize serial communications with the PC

	mfrc522 = setup_rfid_reader();

	Serial.print(">");

	while (!Serial)
		; // Do nothing if no serial port is opened (added for Arduinos based on ATMEGA32U4)
}

void loop()
{
	delay(500);

	// Reset the loop if no new card present on the sensor/reader. This saves the entire process when idle.
	if (!mfrc522.PICC_IsNewCardPresent())
	{
		return;
	}

	// Select one of the cards
	if (!mfrc522.PICC_ReadCardSerial())
	{
		return;
	}
	String uid;
	Serial.print("ID:");
	do
	{
		/* There are Mifare PICCs which have 4 byte or 7 byte UID care if you use 7 byte PICC
			I think we should assume every PICC as they have 4 byte UID
			Until we support 7 byte PICCs */
		for (uint8_t i = 0; i < mfrc522.uid.size; i++)
		{
			Serial.print(String(mfrc522.uid.uidByte[i], DEC));
			if (i < mfrc522.uid.size - 1)
			{
				Serial.print(" ");
			}

			uid += String(mfrc522.uid.uidByte[i], HEX);
		}
	} while ((uid.length() / 2) != mfrc522.uid.size);
	Serial.printf("\r\n>");

#if 0
	Serial.printf("ID:%s\r\n>", uid.c_str());
#endif
}
