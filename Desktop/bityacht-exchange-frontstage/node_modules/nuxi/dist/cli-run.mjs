process._startTime = Date.now();
import('./cli.mjs').then((r) => (r.default || r).main());
