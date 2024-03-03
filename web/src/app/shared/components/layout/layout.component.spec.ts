import { ComponentFixture, TestBed } from '@angular/core/testing';
import { LayoutComponent } from './layout.component';
import { VersionService } from '../../services';
import { SharedSpies } from '../../testings';
import { VersionResponse } from '../../types';
import { of } from 'rxjs';

describe('LayoutComponent', () => {
  let component: LayoutComponent;
  let fixture: ComponentFixture<LayoutComponent>;
  let versionService: jasmine.SpyObj<VersionService>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [LayoutComponent],
      providers: [
        {
          provide: VersionService,
          useValue: SharedSpies.createVersionSpy(),
        },
      ],
    }).compileComponents();
  });

  beforeEach(() => {
    versionService = TestBed.inject(
      VersionService
    ) as jasmine.SpyObj<VersionService>;
    versionService.getVersion.and.returnValue(
      of({
        major: '1',
        build: 'dev',
      } as VersionResponse)
    );

    fixture = TestBed.createComponent(LayoutComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
