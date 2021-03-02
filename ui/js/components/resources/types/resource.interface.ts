export interface RegisterResourceRequest {
  address: string;
  name: string;
  isDefault: boolean;
}

export interface UpdateResourceRequest {
  address: string;
  id: string;
  name: string;
  isDefault: boolean;
}

export interface RemoveResourceRequest {
  id: string;
}

export interface RemoveMemberResourceRequest {
  email: string;
  resourceID: string;
}

export interface BulkAddMembersToResourceRequest {
  emails: string[];
  resourceID: string;
}

export interface ResourceResponse {
  address: string;
  id: string;
  name: string;
  isDefault: boolean;
}

export interface ResourceModalData {
  id: string;
  isEdit: boolean;
  resourceName: string;
  resourceAddress: string;
  isDefault: boolean;
}
