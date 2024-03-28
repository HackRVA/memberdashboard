import { Component, DestroyRef, Inject, OnInit, inject } from '@angular/core';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';
import { MatButtonModule } from '@angular/material/button';
import { MAT_DIALOG_DATA, MatDialogRef } from '@angular/material/dialog';
import { MatIconModule } from '@angular/material/icon';
import { MemberLevelPipe } from '@md-shared/pipes';
import { MemberService } from '@md-shared/services';
import { MemberResponse } from '@md-shared/types';
import { Observable, of, switchMap } from 'rxjs';

@Component({
  selector: 'member-status',
  standalone: true,
  imports: [MatButtonModule, MatIconModule, MemberLevelPipe],
  templateUrl: './member-subscription-status.component.html',
  styleUrl: './member-subscription-status.component.scss',
})
export class MemberSubscriptionStatusComponent implements OnInit {
  private _destroyRef: DestroyRef = inject<DestroyRef>(DestroyRef);
  member: MemberResponse;
  constructor(
    private readonly dialogRef: MatDialogRef<MemberSubscriptionStatusComponent>,
    private readonly memberService: MemberService,
    @Inject(MAT_DIALOG_DATA)
    private readonly dialogData: MemberResponse
  ) {}

  ngOnInit(): void {
    this.fetchAndLoadStatus().subscribe();
  }

  close(): void {
    this.dialogRef.close();
  }

  private fetchAndLoadStatus(): Observable<void> {
    return of(this.dialogData.subscriptionID).pipe(
      switchMap((subscriptionID: string) => {
        if (subscriptionID) {
          return this.memberService
            .checkMemberStatus(this.dialogData.subscriptionID)
            .pipe(
              takeUntilDestroyed(this._destroyRef),
              switchMap((response: MemberResponse) => {
                this.member = response ? response : ({} as MemberResponse);

                return of(null);
              })
            );
        }

        this.member = {
          memberLevel: this.dialogData.memberLevel,
        } as MemberResponse;

        return of(null);
      })
    );
  }
}
