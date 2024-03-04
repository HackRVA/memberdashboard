import { Component, DestroyRef, OnInit, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { Observable, of, switchMap } from 'rxjs';
import { MemberService } from '@md-shared/services';
import { MemberResponse } from '@md-shared/types';
import { MemberLevelPipe } from '@md-shared/pipes';
import { RFIDManagementFactory } from '@md-shared/functions';

@Component({
  selector: 'md-user',
  standalone: true,
  imports: [MatButtonModule, MatDialogModule, MemberLevelPipe],
  templateUrl: './user.component.html',
  styleUrl: './user.component.scss',
})
export class UserComponent implements OnInit {
  private _destroyRef: DestroyRef = inject<DestroyRef>(DestroyRef);
  user: MemberResponse = {} as MemberResponse;
  constructor(
    private readonly memberService: MemberService,
    private readonly dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.fetchAndLoadUser().subscribe();
  }

  async openSelfRFID(): Promise<void> {
    const { MemberRFIDManagementComponent } = await import(
      '@md-shared/components/member-rfid-management'
    );

    this.dialog
      .open(
        MemberRFIDManagementComponent,
        RFIDManagementFactory.createSelfData(this.user.email)
      )
      .afterClosed()
      .pipe(
        switchMap((refresh: boolean) =>
          refresh ? this.fetchAndLoadUser() : of(null)
        )
      )
      .subscribe();
  }

  private fetchAndLoadUser(): Observable<void> {
    return this.memberService.getMemberSelf().pipe(
      takeUntilDestroyed(this._destroyRef),
      switchMap((user: MemberResponse) => {
        if (user) {
          this.user = user;
        }

        return of(null);
      })
    );
  }
}
