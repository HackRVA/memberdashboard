import { Component, DestroyRef, OnInit, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatTableModule } from '@angular/material/table';
import { MatIconModule } from '@angular/material/icon';
import { MatSnackBar } from '@angular/material/snack-bar';
import { MatDialog } from '@angular/material/dialog';
import { MatMenuModule } from '@angular/material/menu';
import { Observable, of, switchMap } from 'rxjs';
import { ResourceResponse } from '@md-shared/types';
import { ResourceService } from '@md-shared/services';
import { ActionBarComponent } from '@md-shared/components/action-bar';
import { ResourceManagementData } from './types';

@Component({
    selector: 'md-resource',
    imports: [
        MatTableModule,
        MatButtonModule,
        MatIconModule,
        MatMenuModule,
        ActionBarComponent,
    ],
    templateUrl: './resource.component.html',
    styleUrl: './resource.component.scss'
})
export class ResourceComponent implements OnInit {
  private _destroyRef: DestroyRef = inject<DestroyRef>(DestroyRef);
  displayedColumns: string[] = ['name', 'address', 'status', 'actions'];
  dataSource: ResourceResponse[] = [];

  constructor(
    private readonly resourceService: ResourceService,
    private readonly snackBar: MatSnackBar,
    private readonly dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.fetchAndLoadResources().subscribe();
  }

  async openResourceManagement(
    resource: ResourceResponse = null
  ): Promise<void> {
    const { ResourceManagementComponent } = await import(
      './components/resource-management'
    );

    this.dialog
      .open(ResourceManagementComponent, {
        autoFocus: false,
        height: '430px',
        width: '320px',
        data: this.generateManagementData(resource),
      })
      .afterClosed()
      .pipe(
        switchMap((refresh: boolean) =>
          refresh ? this.fetchAndLoadResources() : of(null)
        )
      )
      .subscribe();
  }

  async openWarning(
    resource: Pick<ResourceResponse, 'name' | 'id'>
  ): Promise<void> {
    const { ResourceWarningComponent } = await import(
      './components/resource-warning'
    );

    this.dialog
      .open(ResourceWarningComponent, {
        autoFocus: false,
        height: '300px',
        width: '320px',
        data: resource,
      })
      .afterClosed()
      .pipe(
        switchMap((refresh: boolean) =>
          refresh ? this.fetchAndLoadResources() : of(null)
        )
      )
      .subscribe();
  }

  updateACLs(): void {
    this.resourceService.updateACLs().subscribe(() => {
      this.snackBar.open('Successfully update ACL for all resource', '', {
        duration: 3000,
      });
    });
  }

  removeACLs(): void {
    this.resourceService.removeACLs().subscribe(() => {
      this.snackBar.open('Successfully remove ACL for all resource', '', {
        duration: 3000,
      });
    });
  }

  isTimeWitihinActiveRange(lastHeartBeat: string): boolean {
    const lastActiveTime: number = new Date(lastHeartBeat).getTime();
    const currentTime: number = new Date().getTime();

    const thirtyMinsInMS: number = 30 * 60 * 1000;
    const remainingTime: number = currentTime - lastActiveTime;

    return remainingTime > 0 && remainingTime <= thirtyMinsInMS;
  }

  private fetchAndLoadResources(): Observable<void> {
    return this.resourceService.getResources().pipe(
      takeUntilDestroyed(this._destroyRef),
      switchMap((resources: ResourceResponse[]) => {
        if (resources) {
          this.dataSource = resources;
        }

        return of(null);
      })
    );
  }

  private generateManagementData(
    resource: ResourceResponse
  ): ResourceManagementData {
    if (resource) {
      return {
        title: 'Edit a resource',
        name: resource.name,
        address: resource.address,
        id: resource.id,
        isDefault: resource.isDefault,
      };
    }

    return {
      title: 'Register a resource',
    } as ResourceManagementData;
  }
}
