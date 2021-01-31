#include <Arduino.h>
#include <MFRC522.h>

extern void checkID(MFRC522 reader);
extern void granted(uint16_t setDelay);