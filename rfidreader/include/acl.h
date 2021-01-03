/**
 * Our access list is stored in flash memory
 * 
 * acl.cpp will provide functions to read and write to this access list
 * 
 * length of the ACL is stored at address 0
 */
#include <Arduino.h>

/* MAXIMUM_ACL_SIZE is how many IDs that the ACL supports */
#define MAXIMUM_ACL_SIZE 255

/**
 * write_acl
 * we receive a list of IDs
 * note: IDs are 4 bytes long
 */
extern void write_acl(String acl[MAXIMUM_ACL_SIZE], uint8_t arrsize);

/**
 * find_id
 * Determines if the id passed in exists in the ACL
 */
extern bool find_id(String uid);

extern void acl_init();
extern String acl_hash();
