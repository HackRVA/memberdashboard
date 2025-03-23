import { Component, DestroyRef, OnInit, inject } from '@angular/core';
import {
  FormControl,
  FormGroup,
  FormsModule,
  ReactiveFormsModule,
} from '@angular/forms';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatInputModule } from '@angular/material/input';
import { MatIconModule } from '@angular/material/icon';
import { MatTabsModule } from '@angular/material/tabs';
import { MatMenuModule } from '@angular/material/menu';
import { MatFormFieldModule } from '@angular/material/form-field';
import { MatTableModule } from '@angular/material/table';
import { MatSnackBar, MatSnackBarModule } from '@angular/material/snack-bar';
import { MatPaginatorModule, PageEvent } from '@angular/material/paginator';
import { MatDialog } from '@angular/material/dialog';
import {
  BehaviorSubject,
  Observable,
  debounceTime,
  distinctUntilChanged,
  of,
  switchMap,
} from 'rxjs';
import { MemberService, ResourceService } from '@md-shared/services';
import {
  MemberResponse,
  MemberSearchRequest,
  MembersPaginatedResponse,
} from '@md-shared/types';
import { ActionBarComponent } from '@md-shared/components/action-bar';
import { MemberLevelPipe } from '@md-shared/pipes';
import { RFIDManagementFactory } from '@md-shared/functions';
import { MemberResourceManagementData } from './types';

@Component({
  selector: 'md-member',
  standalone: true,
  imports: [
    MatButtonModule,
    MatInputModule,
    MatFormFieldModule,
    MatTableModule,
    MatIconModule,
    MatTabsModule,
    MatMenuModule,
    MatPaginatorModule,
    MatSnackBarModule,
    FormsModule,
    ReactiveFormsModule,
    ActionBarComponent,
    MemberLevelPipe,
  ],
  templateUrl: './member.component.html',
  styleUrl: './member.component.scss',
})
export class MemberComponent implements OnInit {
  private _destroyRef: DestroyRef = inject<DestroyRef>(DestroyRef);
  pageInfo: PageEvent = { pageSize: 10 } as PageEvent;

  private static readonly defaultSearchRequest: MemberSearchRequest = {
    page: 0,
    count: 10,
    active: true,
  } as MemberSearchRequest;
  searchRequest$: BehaviorSubject<MemberSearchRequest> =
    new BehaviorSubject<MemberSearchRequest>(
      MemberComponent.defaultSearchRequest
    );
  searchCtrl: FormControl<string> = new FormControl<string>('');
  memberFormGroup: FormGroup = new FormGroup({
    search: this.searchCtrl,
  });

  displayedColumns: string[] = [
    'name',
    'email',
    'rfid',
    'status',
    'resources',
    'actions',
  ];
  dataSource: MemberResponse[] = [];
  constructor(
    private readonly memberService: MemberService,
    private readonly resourceService: ResourceService,
    private readonly snackBar: MatSnackBar,
    private readonly dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.searchCtrl.valueChanges
      .pipe(debounceTime(500), distinctUntilChanged())
      .subscribe((value: string) => {
        this.searchRequest$.next(
          value.length > 0
            ? ({ search: value } as MemberSearchRequest)
            : MemberComponent.defaultSearchRequest
        );
      });
    this.fetchAndLoadMembers().subscribe();
  }

  pushToResources(member: MemberResponse): void {
    // todo: this should be happen in the API
    for (const resource of member.resources ?? []) {
      this.resourceService
        .bulkAddMembersToResource({
          emails: [member.email],
          resourceID: resource.resourceID,
        })
        .subscribe(() => {
          this.snackBar.open('Success', '', {
            duration: 3000,
          });
        });
    }
  }

  handleTabChange(event: number): void {
    this.pageInfo.pageIndex = 0;
    this.searchRequest$.next({
      active: event === 0 ? true : false,
      page: this.pageInfo.pageIndex,
      count: this.pageInfo.pageSize,
    } as MemberSearchRequest);
  }

  async openAssignMemberRFIDManagement(
    member: MemberResponse = null
  ): Promise<void> {
    const { MemberRFIDManagementComponent } = await import(
      '@md-shared/components/member-rfid-management'
    );

    this.dialog
      .open(
        MemberRFIDManagementComponent,
        member
          ? RFIDManagementFactory.createEditMemberData(
              member.name,
              member.email
            )
          : RFIDManagementFactory.createNewMemberData()
      )
      .afterClosed()
      .pipe(
        switchMap((refresh: boolean) =>
          refresh ? this.fetchAndLoadMembers() : of(null)
        )
      )
      .subscribe();
  }

  async openMemberManagement(member: MemberResponse): Promise<void> {
    const { MemberManagementComponent } = await import(
      './components/member-management'
    );

    this.dialog
      .open(MemberManagementComponent, {
        autoFocus: false,
        width: '320px',
        data: {
          email: member.email,
          name: member.name,
          subscriptionID: member.subscriptionID,
        } as Pick<MemberResponse, 'email' | 'name' | 'subscriptionID'>,
      })
      .afterClosed()
      .pipe(
        switchMap((refresh: boolean) =>
          refresh ? this.fetchAndLoadMembers() : of(null)
        )
      )
      .subscribe();
  }

  async openMemberResourceManagement(
    member: MemberResponse,
    isRemove: boolean = false
  ): Promise<void> {
    const { MemberResourceManagementComponent } = await import(
      './components/member-resource-management'
    );

    this.dialog
      .open(MemberResourceManagementComponent, {
        autoFocus: false,
        width: '320px',
        data: {
          title: isRemove ? 'Remove resource' : 'Add resource',
          email: member.email,
          resources: isRemove ? member.resources : null,
        } as MemberResourceManagementData,
      })
      .afterClosed()
      .pipe(
        switchMap((refresh: boolean) =>
          refresh ? this.fetchAndLoadMembers() : of(null)
        )
      )
      .subscribe();
  }

  async openMemberSubscriptionStatus(member: MemberResponse): Promise<void> {
    const { MemberSubscriptionStatusComponent } = await import(
      './components/member-subscription-status'
    );

    this.dialog.open(MemberSubscriptionStatusComponent, {
      autoFocus: false,
      width: '320px',
      data: member,
    });
  }

  paginate(pageInfo: PageEvent): void {
    const currentRequest: MemberSearchRequest = this.searchRequest$.getValue();
    this.searchRequest$.next({
      active: currentRequest.active,
      page: pageInfo.pageIndex,
      count: pageInfo.pageSize,
    } as MemberSearchRequest);
  }

  private fetchAndLoadMembers(): Observable<void> {
    return this.searchRequest$.pipe(
      switchMap((request: MemberSearchRequest) => {
        if (request) {
          return this.memberService.getMembers(request).pipe(
            takeUntilDestroyed(this._destroyRef),
            switchMap((response: MembersPaginatedResponse) => {
              if (response) {
                this.dataSource = response.members;
                this.pageInfo = Object.assign({}, this.pageInfo, {
                  pageIndex: request.page,
                  length: response.count,
                } as PageEvent);
              }

              return of(null);
            })
          );
        }

        return of(null);
      })
    );
  }
}
