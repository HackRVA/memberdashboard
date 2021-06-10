import { MemberLevel } from "./member.enum";

export interface AssignRFIDRequest {
  email: string;
  rfid: string;
}

export interface MemberResponse {
  id: string;
  name: string;
  email: string;
  rfid: string;
  memberLevel: MemberLevel;
  resources: Array<MemberResource>;
}

export interface MemberResource {
  resourceID: string;
  name: string;
}

export interface CreateMemberRequest {
  email: string;
  rfid: string;
}
