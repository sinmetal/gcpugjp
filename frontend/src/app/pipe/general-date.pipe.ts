import { Pipe, PipeTransform } from '@angular/core';
import * as moment from 'moment';

@Pipe({
  name: 'generalDate'
})
export class GeneralDatePipe implements PipeTransform {

  transform(date: Date): string {
    moment.locale('ja');
    return moment(date).format('YYYY/MM/DD (ddd) HH:mm');
  }

}
