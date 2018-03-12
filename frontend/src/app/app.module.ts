import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';

import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import {
  MatToolbarModule,
  MatProgressBarModule,
  MatButtonModule,
  MatSidenavModule,
  MatListModule,
  MatCardModule,
  MatInputModule
} from '@angular/material';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { HomeComponent } from './home/home.component';
import { AreaEditComponent } from './area-edit/area-edit.component';
import { ForecastComponent } from './forecast/forecast.component';
import { OpenWeatherMapService } from './services/open-weather-map.service';
import { UnixTimeDatePipe } from './pipes/unix-time-date.pipe';
import { LoadingInterceptor } from './loading-interceptor';
import { LoadingService } from './services/loading.service';
import { AreaService } from './services/area.service';
import { ChartsModule } from 'ng2-charts/ng2-charts';
import { PugeventComponent } from './pugevent/pugevent.component';
import { PugeventService } from "./services/pugevent.service";
import { OrganizationService } from "./services/organization.service";
import { GeneralDatePipe } from './pipe/general-date.pipe';
import { OrganizationComponent } from './organization/organization.component';

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    AreaEditComponent,
    ForecastComponent,
    UnixTimeDatePipe,
    PugeventComponent,
    OrganizationComponent,
    GeneralDatePipe,
    OrganizationComponent
  ],
  imports: [
    ChartsModule,
    BrowserModule,
    FormsModule,
    HttpClientModule,
    AppRoutingModule,

    BrowserAnimationsModule,
    MatToolbarModule,
    MatProgressBarModule,
    MatButtonModule,
    MatSidenavModule,
    MatListModule,
    MatCardModule,
    MatInputModule
  ],
  providers: [
    AreaService,
    OpenWeatherMapService,
    PugeventService,
    OrganizationService,
    LoadingService,
    {
      provide: HTTP_INTERCEPTORS,
      useClass: LoadingInterceptor,
      multi: true
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule { }
