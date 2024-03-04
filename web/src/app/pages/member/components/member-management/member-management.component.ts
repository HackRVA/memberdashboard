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
import { MemberService } from '@md-shared/services';
import { MemberResponse } from '@md-shared/types';

@Component({
  selector: 'member-management',
  standalone: true,
  imports: [
    MatButtonModule,
    MatIconModule,
    MatFormFieldModule,
    MatInputModule,
    FormsModule,
    ReactiveFormsModule,
  ],
  templateUrl: './member-management.component.html',
  styleUrl: './member-management.component.scss',
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
    >
  ) {}

  ngOnInit(): void {
    this.handleData(this.dialogData);
  }

  close(): void {
    this.dialogRef.close();
  }

  submit(): void {
    this.memberService
      .updateMemberByEmail(
        this.dialogData.email,
        this.memberManagementGroup.value
      )
      .subscribe(() => {
        this.close();
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
