import { ComponentFixture, TestBed } from '@angular/core/testing';
import { MinionComponent } from './minion.component';

describe('MinionComponent', () => {
  let component: MinionComponent;
  let fixture: ComponentFixture<MinionComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [MinionComponent],
    }).compileComponents();

    fixture = TestBed.createComponent(MinionComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
