import { TestBed } from '@angular/core/testing';
import { AppComponent } from './app.component';
import { AuthService } from './shared/services';
import { SharedSpies } from './shared/testings';

describe('AppComponent', () => {
  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppComponent],
      providers: [
        {
          provide: AuthService,
          useValue: SharedSpies.createAuthServiceSpy(),
        },
      ],
    }).compileComponents();
  });

  it('should create the app', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.componentInstance;
    expect(app).toBeTruthy();
  });
});
