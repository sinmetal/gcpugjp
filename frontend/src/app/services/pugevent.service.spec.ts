import { TestBed, inject } from '@angular/core/testing';

import { PugeventService } from './pugevent.service';

describe('PugeventService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [PugeventService]
    });
  });

  it('should be created', inject([PugeventService], (service: PugeventService) => {
    expect(service).toBeTruthy();
  }));
});
