export type MemberSearchRequest = {
  page: number;
  count: number;
  active: boolean;
  search: string;
};

export enum MemberLevel {
  Inactive = 1,
  Credited = 2,
  Classic = 3,
  Standard = 4,
  Premium = 5,
}

export type MemberResource = {
  resourceID: string;
  name: string;
};

export type MemberResponse = {
  id: string;
  name: string;
  email: string;
  rfid: string;
  subscriptionID: string;
  memberLevel: MemberLevel;
  resources: MemberResource[];
};

export type AssignRFIDRequest = {
  email: string;
  rfid: string;
};

export type CreateMemberRequest = AssignRFIDRequest & { name: string };

export type UpdateMemberRequest = {
  fullName: string;
  subscriptionID: string;
};

export type AckResponse = {
  ack: boolean;
};
