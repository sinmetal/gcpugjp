import { Pipe, PipeTransform } from '@angular/core';
import { DatePipe } from '@angular/common';

@Pipe({
  name: 'unixTimeDate'
})
export class UnixTimeDatePipe implements PipeTransform {
  private datePipe: DatePipe = new DatePipe('ja-JP');
  transform(value: any, pattern?: string): any {
    // UnixTimeから戻しつつDatePipeを実行
    return this.datePipe.transform(new Date(value * 1000), pattern);
  }
}