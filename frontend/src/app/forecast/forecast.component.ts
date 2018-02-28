import { Component, OnInit } from '@angular/core';

import { ActivatedRoute } from '@angular/router';
import { OpenWeatherMapService } from '../services/open-weather-map.service';
import { OpenWeatherMap } from '../shared/models/open-weather-map';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/switchMap';

import {
  transition,
  trigger,
  query,
  style,
  stagger,
  useAnimation
} from '@angular/animations';
import { slideFadeIn } from '../app.animations';

import { ForecastChartData } from '../shared/models/forecast-chart-data';
import { UnixTimeDatePipe } from '../pipes/unix-time-date.pipe';

@Component({
  selector: 'app-forecast',
  templateUrl: './forecast.component.html',
  styleUrls: ['./forecast.component.scss'],
  animations: [
    trigger('slideFade', [
      transition(':enter', [
        query('mat-list-item', [
          style({
            opacity: 0
          }),
          stagger(100, [
            useAnimation(slideFadeIn)
          ])
        ])
      ])
    ])
  ]
})
export class ForecastComponent implements OnInit {
  public currentWeatherObservable: Observable<OpenWeatherMap.Current>;
  public forecastObservable: Observable<OpenWeatherMap.Forecast>;

  public lineChartLabels: string[] = [];
  public lineChartData: ForecastChartData[] = [];

  private unixTimeDatePipe: UnixTimeDatePipe = new UnixTimeDatePipe();

  constructor(
    private route: ActivatedRoute,
    private openWeatherMapService: OpenWeatherMapService
  ) { }

  ngOnInit() {
    // 現在の天気
    this.currentWeatherObservable = this.route.params.switchMap(param => {
      return this.openWeatherMapService.current(param['city']);
    });

    // 1週間の天気予報を取得
    this.forecastObservable = this.route.params.switchMap(param => {
      return this.openWeatherMapService.forecast(param['city'])
        .map((res: OpenWeatherMap.Forecast) => {
          this.lineChartLabels = [];

          const maxTemp: ForecastChartData = {
            data: [],
            label: '最高気温'
          };
          const minTemp: ForecastChartData = {
            data: [],
            label: '最低気温'
          };

          res.list.forEach((weather: OpenWeatherMap.DailyWeather) => {
            const day = this.unixTimeDatePipe.transform(weather.dt, 'MM/dd');
            this.lineChartLabels.push(day);
            maxTemp.data.push(weather.temp.max);
            minTemp.data.push(weather.temp.min);
          });

          this.lineChartData = [
            maxTemp,
            minTemp
          ];
          return res;
        })
    });
  }
}