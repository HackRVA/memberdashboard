import { MatDialogConfig } from '@angular/material/dialog';
import { RFIDManagementData, RFIDManagementType } from '../types';

export class RFIDManagementFactory {
  public static createSelfData(
    name: string,
    email: string
  ): MatDialogConfig<RFIDManagementData> {
    return {
      autoFocus: false,
      width: '320px',
      data: {
        name: name,
        email: email,
        title: 'Assign RFID',
        shouldDisable: true,
        type: RFIDManagementType.Self,
      },
    };
  }

  public static createNewMemberData(): MatDialogConfig<RFIDManagementData> {
    return {
      autoFocus: false,
      width: '320px',
      data: {
        title: 'Assign new member',
        shouldDisable: false,
        type: RFIDManagementType.New,
      } as RFIDManagementData,
    };
  }

  public static createEditMemberData(
    name: string,
    email: string
  ): MatDialogConfig<RFIDManagementData> {
    return {
      autoFocus: false,
      width: '320px',
      data: {
        name: name,
        email: email,
        title: 'Assign member RFID',
        shouldDisable: true,
        type: RFIDManagementType.Edit,
      } as RFIDManagementData,
    };
  }
}
