import { TestBed } from '@angular/core/testing';
import { MemberService } from './member.service';
import { HttpClient } from '@angular/common/http';
import { SharedSpies } from '../testings';

describe('MemberService', () => {
  let service: MemberService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        { provide: HttpClient, useValue: SharedSpies.createHttpClientSpy() },
      ],
    });
    service = TestBed.inject(MemberService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
