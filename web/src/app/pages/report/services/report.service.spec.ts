import { TestBed } from '@angular/core/testing';
import { ReportService } from './report.service';
import { HttpClient } from '@angular/common/http';
import { SharedSpies } from '@md-shared/testings';

describe('ReportService', () => {
  let service: ReportService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        ReportService,
        { provide: HttpClient, useValue: SharedSpies.createHttpClientSpy() },
      ],
    });
    service = TestBed.inject(ReportService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
