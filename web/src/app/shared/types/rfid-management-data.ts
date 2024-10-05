export type RFIDManagementData = {
  name: string;
  title: string;
  email: string;
  shouldDisable: boolean;
  type: RFIDManagementType;
};

export enum RFIDManagementType {
  Self,
  New,
  Edit,
}
