import { TestBed, inject } from '@angular/core/testing';

import { OpenWeatherMapService } from './open-weather-map.service';

describe('OpenWeatherMapService', () => {
  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [OpenWeatherMapService]
    });
  });

  it('should be created', inject([OpenWeatherMapService], (service: OpenWeatherMapService) => {
    expect(service).toBeTruthy();
  }));
});
