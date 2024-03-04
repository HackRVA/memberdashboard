import { Component, Inject, OnInit } from '@angular/core';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatIconModule } from '@angular/material/icon';
import { ResourceService } from '@md-shared/services';
import { ResourceResponse } from '@md-shared/types';

@Component({
  selector: 'resource-warning',
  templateUrl: './resource-warning.component.html',
  standalone: true,
  imports: [MatButtonModule, MatIconModule, MatSnackBarModule],
  styleUrl: './resource-warning.component.scss',
})
export class ResourceWarningComponent implements OnInit {
  name: string = '';
  constructor(
    private readonly resourceService: ResourceService,
    private readonly dialogRef: MatDialogRef<ResourceWarningComponent>,
    private readonly snackBar: MatSnackBar,
    @Inject(MAT_DIALOG_DATA)
    private readonly dialogData: Pick<ResourceResponse, 'name' | 'id'>
  ) {}

  ngOnInit(): void {
    this.name = this.dialogData.name;
  }

  close(): void {
    this.dialogRef.close();
  }

  submit(): void {
    this.resourceService
      .deleteResource({ id: this.dialogData.id })
      .subscribe(() => {
        this.snackBar.open(
          'Success, but I hope you know what you are doing :3',
          '',
          { duration: 3000 }
        );
        this.close();
      });
  }
}
