import { UnixTimeDatePipe } from './unix-time-date.pipe';

describe('UnixTimeDatePipe', () => {
  it('create an instance', () => {
    const pipe = new UnixTimeDatePipe();
    expect(pipe).toBeTruthy();
  });
});
