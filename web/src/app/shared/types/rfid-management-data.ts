export type RFIDManagementData = {
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
