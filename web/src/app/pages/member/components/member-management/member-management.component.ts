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
import { MemberService } from '@md-shared/services';
import { MemberResponse } from '@md-shared/types';

@Component({
    selector: 'member-management',
    imports: [
        MatButtonModule,
        MatIconModule,
        MatFormFieldModule,
        MatInputModule,
        MatSnackBarModule,
        FormsModule,
        ReactiveFormsModule,
        CommonModule,
    ],
    templateUrl: './member-management.component.html',
    styleUrl: './member-management.component.scss'
})
export class MemberManagementComponent implements OnInit {
  memberManagementGroup: FormGroup = new FormGroup({
    fullName: new FormControl<string>(null, [Validators.required]),
    subscriptionID: new FormControl<string>(null, [Validators.required]),
  });

  constructor(
    private readonly dialogRef: MatDialogRef<MemberManagementComponent>,
    private readonly memberService: MemberService,
    @Inject(MAT_DIALOG_DATA)
    private readonly dialogData: Pick<
      MemberResponse,
      'name' | 'subscriptionID' | 'email'
    >,
    private readonly snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.handleData(this.dialogData);
  }

  close(refresh: boolean = false): void {
    this.dialogRef.close(refresh);
  }

  submit(): void {
    this.memberService
      .updateMemberByEmail(
        this.dialogData.email,
        this.memberManagementGroup.value
      )
      .subscribe({
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

  private handleData(
    data: Pick<MemberResponse, 'name' | 'subscriptionID' | 'email'>
  ) {
    this.memberManagementGroup.get('fullName').setValue(data.name);
    this.memberManagementGroup
      .get('subscriptionID')
      .setValue(data.subscriptionID);
  }
}
