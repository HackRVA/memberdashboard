import { MemberLevel } from '../custom/member-level';

export interface MemberResource {
  resourceID: string;
  name: string;
}

export interface Member {
  id: string;
  name: string;
  email: string;
  rfid: string;
  subscriptionID: string;
  memberLevel: MemberLevel;
  resources: Array<MemberResource>;
}
