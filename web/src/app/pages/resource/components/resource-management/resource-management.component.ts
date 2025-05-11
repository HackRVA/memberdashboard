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
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatIconModule } from '@angular/material/icon';
import { MatInputModule } from '@angular/material/input';
import { MatCheckboxModule } from '@angular/material/checkbox';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { Observable } from 'rxjs';
import { ResourceService } from '@md-shared/services';
import { ResourceManagementData } from '../../types';

@Component({
    selector: 'resource-management',
    imports: [
        CommonModule,
        MatIconModule,
        MatButtonModule,
        MatFormFieldModule,
        MatInputModule,
        MatCheckboxModule,
        MatSnackBarModule,
        FormsModule,
        ReactiveFormsModule,
    ],
    templateUrl: './resource-management.component.html',
    styleUrl: './resource-management.component.scss'
})
export class ResourceManagementComponent implements OnInit {
  title: string = '';
  resourceManagementGroup: FormGroup = new FormGroup({
    name: new FormControl<string>(null, [Validators.required]),
    address: new FormControl<string>(null, [Validators.required]),
    isDefault: new FormControl<boolean>(false),
  });

  constructor(
    private readonly resourceService: ResourceService,
    private readonly dialogRef: MatDialogRef<ResourceManagementComponent>,
    @Inject(MAT_DIALOG_DATA)
    private readonly dialogData: ResourceManagementData,
    private readonly snackBar: MatSnackBar
  ) {}

  ngOnInit(): void {
    this.handleData(this.dialogData);
  }

  close(refresh: boolean = false): void {
    this.dialogRef.close(refresh);
  }

  submit(): void {
    let resourceObs$: Observable<void> = null;

    if (this.dialogData.id) {
      resourceObs$ = this.resourceService.updateResource({
        id: this.dialogData.id,
        ...this.resourceManagementGroup.value,
      });
    } else {
      resourceObs$ = this.resourceService.registerResource(
        this.resourceManagementGroup.value
      );
    }

    resourceObs$.subscribe({
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

  private handleData(data: ResourceManagementData): void {
    this.title = data.title;

    if (data) {
      this.resourceManagementGroup.get('name').setValue(data.name);
      this.resourceManagementGroup.get('address').setValue(data.address);
      this.resourceManagementGroup.get('isDefault').setValue(data.isDefault);
    }
  }
}
