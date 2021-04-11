#include <Arduino.h>
#include <MFRC522.h>

#include "rfid.h"
#include "acl.h"
#include "wifi.h"
#include "relay.h"

#define TESTING 0
#define DEBUG 0

MFRC522 mfrc522;

void setup()
{
	Serial.begin(9600); // Initialize serial communications with the PC
	mfrc522 = setup_rfid_reader();
	acl_init();
	setup_wifi();
	pinMode(RELAY_PIN, OUTPUT);

#if TESTING
	unsigned long acl[MAXIMUM_ACL_SIZE] = {};
	write_acl(acl, 0);
#endif

	while (!Serial)
		; // Do nothing if no serial port is opened (added for Arduinos based on ATMEGA32U4)
}

void loop()
{
	wifi_loop();
	read_rfid(mfrc522);
}
