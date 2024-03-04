import { Component, DestroyRef, Inject, OnInit, inject } from '@angular/core';
import {
  FormControl,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
  Validators,
} from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { CommonModule } from '@angular/common';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatInputModule } from '@angular/material/input';
import { MatSelectModule } from '@angular/material/select';
import { Observable, of, switchMap } from 'rxjs';
import { ResourceService } from '@md-shared/services';
import {
  BulkAddMembersToResourceRequest,
  MemberResource,
  ResourceResponse,
} from '@md-shared/types';
import { MemberResourceManagementData } from '../../types';

@Component({
  selector: 'member-resource-management',
  standalone: true,
  imports: [
    CommonModule,
    MatSelectModule,
    MatButtonModule,
    MatFormFieldModule,
    MatInputModule,
    FormsModule,
    ReactiveFormsModule,
  ],
  templateUrl: './member-resource-management.component.html',
  styleUrl: './member-resource-management.component.scss',
})
export class MemberResourceManagementComponent implements OnInit {
  private _destroyRef: DestroyRef = inject<DestroyRef>(DestroyRef);
  title: string = '';
  memberResourceManagementGroup: FormGroup = new FormGroup({
    email: new FormControl<string>({ value: '', disabled: true }, [
      Validators.required,
    ]),
    resourceID: new FormControl<string>(null, [Validators.required]),
  });
  resourcesObs$: Observable<MemberResource[]>;

  constructor(
    private readonly dialogRef: MatDialogRef<MemberResourceManagementComponent>,
    private readonly resourceService: ResourceService,
    @Inject(MAT_DIALOG_DATA)
    private readonly dialogData: MemberResourceManagementData
  ) {}

  ngOnInit(): void {
    this.handleData(this.dialogData);
  }

  close(refresh: boolean = false): void {
    this.dialogRef.close(refresh);
  }

  submit(): void {
    let resourceService$: Observable<void> = null;

    if (this.dialogData.resources) {
      resourceService$ = this.resourceService.removeMemberFromResource({
        email: this.memberResourceManagementGroup.get('email').value,
        resourceID: this.memberResourceManagementGroup.get('resourceID').value,
      });
    } else {
      resourceService$ = this.resourceService.bulkAddMembersToResource({
        emails: [this.memberResourceManagementGroup.get('email').value],
        resourceID: this.memberResourceManagementGroup.get('resourceID').value,
      } as BulkAddMembersToResourceRequest);
    }

    resourceService$.subscribe(() => {
      this.close(true);
    });
  }

  private handleData(data: MemberResourceManagementData): void {
    this.title = data.title;
    this.memberResourceManagementGroup.get('email').setValue(data.email);
    this.resourcesObs$ = this.fetchResources(data);
  }

  private fetchResources(
    data: MemberResourceManagementData
  ): Observable<MemberResource[]> {
    if (data.resources) {
      return of(data.resources);
    }
    return this.resourceService.getResources().pipe(
      takeUntilDestroyed(this._destroyRef),
      switchMap((response: ResourceResponse[]) =>
        of(
          response.map((x: ResourceResponse) => {
            return { name: x.name, resourceID: x.id } as MemberResource;
          })
        )
      )
    );
  }
}
