import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { TopNavComponent } from './top-nav.component';
import { AuthService } from '../../services';
import { SharedSpies } from '../../testings';
import { BehaviorSubject } from 'rxjs';
import { AuthUser } from '../../types';

describe('TopNavComponent', () => {
  let component: TopNavComponent;
  let fixture: ComponentFixture<TopNavComponent>;
  let authService: jasmine.SpyObj<AuthService>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [TopNavComponent, RouterTestingModule],
      providers: [
        {
          provide: AuthService,
          useValue: SharedSpies.createAuthServiceSpy(),
        },
      ],
    }).compileComponents();
  });

  beforeEach(() => {
    authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
    authService.user$ = new BehaviorSubject({} as AuthUser);

    fixture = TestBed.createComponent(TopNavComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
