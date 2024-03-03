import { TestBed } from '@angular/core/testing';
import { ResourceService } from './resource.service';
import { HttpClient } from '@angular/common/http';
import { SharedSpies } from '../testings';

describe('ResourceService', () => {
  let service: ResourceService;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        { provide: HttpClient, useValue: SharedSpies.createHttpClientSpy() },
      ],
    });
    service = TestBed.inject(ResourceService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });
});
