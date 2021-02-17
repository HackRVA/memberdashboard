import { MemberLevel } from "./member.enum";

export interface AssignRFIDRequest {
  email: string;
  rfid: string;
}

export interface MemberResponse {
  id: string;
  name: string;
  email: string;
  memberLevel: MemberLevel;
  resources: Array<MemberResource>;
}

export interface MemberResource {
  resourceID: string;
  name: string;
}
