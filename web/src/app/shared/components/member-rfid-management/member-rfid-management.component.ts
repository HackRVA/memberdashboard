import { CommonModule } from '@angular/common';
import { Component, Inject, OnInit } from '@angular/core';
import {
  FormControl,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MemberService } from '../../services';
import {
  AssignRFIDRequest,
  RFIDManagementData,
  RFIDManagementType,
} from '../../types';
import { Observable, of } from 'rxjs';

@Component({
  selector: 'member-rfid-management',
  standalone: true,
  imports: [
    CommonModule,
    MatIconModule,
    MatButtonModule,
    MatFormFieldModule,
    MatInputModule,
    MatSnackBarModule,
    FormsModule,
    ReactiveFormsModule,
  ],
  templateUrl: './member-rfid-management.component.html',
  styleUrl: './member-rfid-management.component.scss',
})
export class MemberRFIDManagementComponent implements OnInit {
  title: string = '';
  rfidType = RFIDManagementType;
  memberRFIDType: RFIDManagementType;

  rfidManagementGroup: FormGroup = new FormGroup({
    name: new FormControl<string>(null, [Validators.required]),
    email: new FormControl<string>(null, [
      Validators.required,
      Validators.email,
    ]),
    rfid: new FormControl<number>(null, [Validators.required]),
  });

  constructor(
    private readonly dialogRef: MatDialogRef<MemberRFIDManagementComponent>,
    private readonly memberService: MemberService,
    @Inject(MAT_DIALOG_DATA)
    private readonly dialogData: RFIDManagementData,
    private readonly snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.handleData(this.dialogData);
  }

  close(refresh: boolean = false): void {
    this.dialogRef.close(refresh);
  }

  submit(): void {
    let memberObs$: Observable<void> = null;
    const request: AssignRFIDRequest = {
      email: this.rfidManagementGroup.get('email').value,
      rfid: this.rfidManagementGroup.get('rfid').value.toString(),
    };

    switch (this.memberRFIDType) {
      case RFIDManagementType.Self:
        memberObs$ = this.memberService.assignRFIDToSelf(request);
        break;
      case RFIDManagementType.New:
        memberObs$ = this.memberService.assignNewMemberRFID({
          ...request,
          name: this.rfidManagementGroup.get('name').value,
        });
        break;
      case RFIDManagementType.Edit:
        memberObs$ = this.memberService.assignRFID(request);
        break;
      default:
        memberObs$ = of(null);
        break;
    }
    memberObs$.subscribe({
      next: () => {
        this.snackBar.open('Success', '', { duration: 3000 });
        this.close(true);
      },
      error: () => {
        this.snackBar.open('Hrmmm, it failed', '', { duration: 3000 });
        this.close(false);
      },
    });
  }

  private handleData(data: RFIDManagementData): void {
    if (data.name) {
      this.rfidManagementGroup.get('name').setValue(data.name);
    }

    if (data.email) {
      this.rfidManagementGroup.get('email').setValue(data.email);
    }

    if (data.shouldDisable) {
      this.rfidManagementGroup.get('name').disable();
      this.rfidManagementGroup.get('email').disable();
    }

    this.title = data.title;
    this.memberRFIDType = data.type;
  }
}
