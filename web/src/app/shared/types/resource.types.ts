export type RegisterResourceRequest = {
  address: string;
  name: string;
  isDefault: boolean;
};

export type RemoveResourceRequest = {
  id: string;
};

export type ResourceResponse = {
  address: string;
  id: string;
  name: string;
  isDefault: boolean;
  lastHeartBeat: string;
};

export type UpdateResourceRequest = Pick<
  ResourceResponse,
  'address' | 'id' | 'name' | 'isDefault'
>;

export type BulkAddMembersToResourceRequest = {
  emails: string[];
  resourceID: string;
};

export type RemoveMemberResourceRequest = {
  email: string;
  resourceID: string;
};
