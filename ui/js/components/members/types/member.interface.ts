import { MemberLevel } from "./member.enum";
import { ResourceResponse } from "../../resources/types";

export interface RemoveMemberResourceModalData {
  email: string;
  memberResources: MemberResource[];
  handleResourceChange: Function;
  handleSubmitRemoveMemberResource: Function;
  emptyFormValuesOnClosed: Function;
}

export interface AddMemberResourceModalData {
  email: string;
  resources: ResourceResponse[];
  handleResourceChange: Function;
  handleSubmitAddMemberResource: Function;
  emptyFormValuesOnClosed: Function;
}

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
