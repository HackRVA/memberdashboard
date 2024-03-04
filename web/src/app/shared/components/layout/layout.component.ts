import { Component, DestroyRef, OnInit, inject } from '@angular/core';
import { VersionService } from '../../services/version.service';
import { VersionResponse } from '../../types/version-response';
import { Observable, of, switchMap } from 'rxjs';
import { takeUntilDestroyed } from '@angular/core/rxjs-interop';

@Component({
  selector: 'md-layout',
  standalone: true,
  imports: [],
  templateUrl: './layout.component.html',
  styleUrl: './layout.component.scss',
})
export class LayoutComponent implements OnInit {
  private _destroyRef: DestroyRef = inject<DestroyRef>(DestroyRef);
  version: VersionResponse;
  constructor(private readonly versionService: VersionService) {}

  ngOnInit(): void {
    this.fetchAndLoadVersion().subscribe();
  }

  private fetchAndLoadVersion(): Observable<void> {
    return this.versionService.getVersion().pipe(
      takeUntilDestroyed(this._destroyRef),
      switchMap((response: VersionResponse) => {
        if (response) {
          this.version = response;
        }

        return of(null);
      })
    );
  }
}
