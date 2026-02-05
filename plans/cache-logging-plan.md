# add logging to internal/explore/cache.go

### Plan

- Add logging before
- Add logging after
- Add logging during if it outputs something

ALL functions need logging.

---

## gcsCache.path (line 54)

```go
func (g *gcsCache) path(key string) string {
	log.Printf("[GCS] path: key=%s", key)
	result := path.Join("soci", strings.Replace(key, ":", "-", 1), "toc.json.gz")
	log.Printf("[GCS] path: result=%s", result)
	return result
}
```

---

## gcsCache.treePath (line 58)

```go
func (g *gcsCache) treePath(key string) string {
	log.Printf("[GCS] treePath: key=%s", key)
	result := path.Join("soci", strings.Replace(key, ":", "-", 1)) + ".tar.gz"
	log.Printf("[GCS] treePath: result=%s", result)
	return result
}
```

---

## gcsCache.object (line 62)

```go
func (g *gcsCache) object(key string) *storage.ObjectHandle {
	log.Printf("[GCS] object: key=%s", key)
	return g.bucket.Object(g.path(key))
}
```

---

## gcsCache.Get (line 67)

```go
func (g *gcsCache) Get(ctx context.Context, key string) (*soci.TOC, error) {
	log.Printf("[GCS] Get: START key=%s", key)
	if debug {
		start := time.Now()
		defer func() {
			log.Printf("bucket.Get(%q) (%s)", key, time.Since(start))
		}()
	}
	log.Printf("[GCS] Get: BEFORE NewReader")
	rc, err := g.object(key).NewReader(ctx)
	log.Printf("[GCS] Get: AFTER NewReader err=%v", err)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	log.Printf("[GCS] Get: BEFORE gzip.NewReader")
	zr, err := gzip.NewReader(rc)
	log.Printf("[GCS] Get: AFTER gzip.NewReader err=%v", err)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	toc := &soci.TOC{}
	log.Printf("[GCS] Get: BEFORE json.Decode")
	if err := json.NewDecoder(zr).Decode(toc); err != nil {
		log.Printf("[GCS] Get: AFTER json.Decode err=%v", err)
		return nil, err
	}
	log.Printf("[GCS] Get: AFTER json.Decode err=nil")
	return toc, nil
}
```

---

## gcsCache.Put (line 92)

```go
func (g *gcsCache) Put(ctx context.Context, key string, toc *soci.TOC) error {
	log.Printf("[GCS] Put: START key=%s", key)
	if debug {
		start := time.Now()
		defer func() {
			log.Printf("bucket.Put(%q) (%s)", key, time.Since(start))
		}()
	}
	log.Printf("[GCS] Put: BEFORE NewWriter")
	w := g.object(key).NewWriter(ctx)
	log.Printf("[GCS] Put: AFTER NewWriter")

	log.Printf("[GCS] Put: BEFORE gzip.NewWriterLevel")
	zw, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	log.Printf("[GCS] Put: AFTER gzip.NewWriterLevel err=%v", err)
	if err != nil {
		return err
	}
	log.Printf("[GCS] Put: BEFORE json.Encode")
	if err := json.NewEncoder(zw).Encode(toc); err != nil {
		log.Printf("[GCS] Put: AFTER json.Encode err=%v", err)
		zw.Close()
		return err
	}
	log.Printf("[GCS] Put: AFTER json.Encode err=nil")
	log.Printf("[GCS] Put: BEFORE zw.Close")
	if err := zw.Close(); err != nil {
		log.Printf("[GCS] Put: AFTER zw.Close err=%v", err)
		return err
	}
	log.Printf("[GCS] Put: AFTER zw.Close err=nil")
	log.Printf("[GCS] Put: BEFORE w.Close")
	err = w.Close()
	log.Printf("[GCS] Put: AFTER w.Close err=%v", err)
	return err
}
```

---

## gcsCache.Writer (line 120)

```go
func (g *gcsCache) Writer(ctx context.Context, key string) (io.WriteCloser, error) {
	log.Printf("[GCS] Writer: key=%s", key)
	return g.bucket.Object(g.treePath(key)).NewWriter(ctx), nil
}
```

---

## gcsCache.Reader (line 124)

```go
func (g *gcsCache) Reader(ctx context.Context, key string) (io.ReadCloser, error) {
	log.Printf("[GCS] Reader: key=%s", key)
	return g.bucket.Object(g.treePath(key)).NewReader(ctx)
}
```

---

## gcsCache.RangeReader (line 128)

```go
func (g *gcsCache) RangeReader(ctx context.Context, key string, offset, length int64) (io.ReadCloser, error) {
	log.Printf("[GCS] RangeReader: key=%s offset=%d length=%d", key, offset, length)
	return g.bucket.Object(g.treePath(key)).NewRangeReader(ctx, offset, length)
}
```

---

## gcsCache.Size (line 132)

```go
func (g *gcsCache) Size(ctx context.Context, key string) (int64, error) {
	log.Printf("[GCS] Size: START key=%s", key)
	if debug {
		start := time.Now()
		defer func() {
			log.Printf("bucket.Size(%q) (%s)", key, time.Since(start))
		}()
	}
	log.Printf("[GCS] Size: BEFORE Attrs")
	attrs, err := g.bucket.Object(g.treePath(key)).Attrs(ctx)
	log.Printf("[GCS] Size: AFTER Attrs err=%v", err)
	if err != nil {
		return -1, err
	}
	log.Printf("[GCS] Size: size=%d", attrs.Size)
	return attrs.Size, nil
}
```

---

## gcsCache.Delete (line 146)

```go
func (g *gcsCache) Delete(ctx context.Context, key string) error {
	log.Printf("[GCS] Delete: key=%s", key)
	log.Printf("[GCS] Delete: BEFORE Delete")
	err := g.bucket.Object(g.treePath(key)).Delete(ctx)
	log.Printf("[GCS] Delete: AFTER Delete err=%v", err)
	return err
}
```

---

## dirCache.file (line 154)

```go
func (d *dirCache) file(key string) string {
	log.Printf("[CACHE] file: key=%s dir=%s", key, d.dir)
	sanitized := strings.Replace(key, ":", "-", 1)
	sanitized = strings.ReplaceAll(sanitized, "/", "-")
	result := filepath.Join(d.dir, sanitized)
	log.Printf("[CACHE] file: result=%s", result)
	return result
}
```

---

## dirCache.Get (line 160)

```go
func (d *dirCache) Get(ctx context.Context, key string) (*soci.TOC, error) {
	path := d.file(key) + ".toc.json.gz"
	log.Printf("[CACHE] Get: BEFORE os.Open path=%s", path)
	f, err := os.Open(path)
	log.Printf("[CACHE] Get: AFTER os.Open err=%v", err)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	log.Printf("[CACHE] Get: BEFORE gzip.NewReader")
	zr, err := gzip.NewReader(f)
	log.Printf("[CACHE] Get: AFTER gzip.NewReader err=%v", err)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	toc := &soci.TOC{}
	log.Printf("[CACHE] Get: BEFORE json.Decode")
	if err := json.NewDecoder(zr).Decode(toc); err != nil {
		log.Printf("[CACHE] Get: AFTER json.Decode err=%v", err)
		return nil, err
	}
	log.Printf("[CACHE] Get: AFTER json.Decode err=nil")
	return toc, nil
}
```

---

## dirCache.Put (line 178)

```go
func (d *dirCache) Put(ctx context.Context, key string, toc *soci.TOC) error {
	path := d.file(key) + ".toc.json.gz"
	log.Printf("[CACHE] Put: BEFORE os.OpenFile path=%s", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0755)
	log.Printf("[CACHE] Put: AFTER os.OpenFile err=%v", err)
	if err != nil {
		return err
	}
	defer f.Close()
	log.Printf("[CACHE] Put: BEFORE gzip.NewWriterLevel")
	zw, err := gzip.NewWriterLevel(f, gzip.BestSpeed)
	log.Printf("[CACHE] Put: AFTER gzip.NewWriterLevel err=%v", err)
	if err != nil {
		return err
	}
	defer zw.Close()
	log.Printf("[CACHE] Put: BEFORE json.Encode")
	err = json.NewEncoder(zw).Encode(toc)
	log.Printf("[CACHE] Put: AFTER json.Encode err=%v", err)
	return err
}
```

---

## dirCache.Writer (line 192)

```go
func (d *dirCache) Writer(ctx context.Context, key string) (io.WriteCloser, error) {
	log.Printf("[CACHE] Writer: START dir=%s key=%s", d.dir, key)
	pattern := strings.Replace(key, ":", "-", 1)
	pattern = strings.ReplaceAll(pattern, "/", "-")
	log.Printf("[CACHE] Writer: pattern=%s", pattern)
	log.Printf("[CACHE] Writer: BEFORE os.CreateTemp dir=%s pattern=%s", d.dir, pattern)
	tmp, err := os.CreateTemp(d.dir, pattern)
	log.Printf("[CACHE] Writer: AFTER os.CreateTemp err=%v", err)
	if err != nil {
		return nil, err
	}
	log.Printf("[CACHE] Writer: tmp.Name()=%s", tmp.Name())
	dst := d.file(key) + ".tar.gz"
	log.Printf("[CACHE] Writer: dst=%s", dst)
	return &dirWriter{
		dst: dst,
		f:   tmp,
	}, nil
}
```

---

## dirCache.Reader (line 206)

```go
func (d *dirCache) Reader(ctx context.Context, key string) (io.ReadCloser, error) {
	path := d.file(key) + ".tar.gz"
	log.Printf("[CACHE] Reader: BEFORE os.Open path=%s", path)
	f, err := os.Open(path)
	log.Printf("[CACHE] Reader: AFTER os.Open err=%v", err)
	return f, err
}
```

---

## dirCache.RangeReader (line 211)

```go
func (d *dirCache) RangeReader(ctx context.Context, key string, offset, length int64) (io.ReadCloser, error) {
	log.Printf("[CACHE] RangeReader: key=%s offset=%d length=%d", key, offset, length)
	if offset == 0 && length == -1 {
		return d.Reader(ctx, key)
	}
	path := d.file(key) + ".tar.gz"
	log.Printf("[CACHE] RangeReader: BEFORE os.Open path=%s", path)
	f, err := os.Open(path)
	log.Printf("[CACHE] RangeReader: AFTER os.Open err=%v", err)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(io.NewSectionReader(f, offset, length)), nil
}
```

---

## dirCache.Size (line 222)

```go
func (d *dirCache) Size(ctx context.Context, key string) (int64, error) {
	path := d.file(key) + ".tar.gz"
	log.Printf("[CACHE] Size: BEFORE os.Stat path=%s", path)
	stat, err := os.Stat(path)
	log.Printf("[CACHE] Size: AFTER os.Stat err=%v", err)
	if err != nil {
		return -1, err
	}
	log.Printf("[CACHE] Size: size=%d", stat.Size())
	return stat.Size(), nil
}
```

---

## dirCache.Delete (line 230)

```go
func (d *dirCache) Delete(ctx context.Context, key string) error {
	tocPath := d.file(key) + ".toc.json.gz"
	tarPath := d.file(key) + ".tar.gz"
	log.Printf("[CACHE] Delete: BEFORE os.Remove tocPath=%s", tocPath)
	err1 := os.Remove(tocPath)
	log.Printf("[CACHE] Delete: AFTER os.Remove tocPath err=%v", err1)
	log.Printf("[CACHE] Delete: BEFORE os.Remove tarPath=%s", tarPath)
	err2 := os.Remove(tarPath)
	log.Printf("[CACHE] Delete: AFTER os.Remove tarPath err=%v", err2)
	return err2
}
```

---

## dirWriter.Write (line 242)

```go
func (d *dirWriter) Write(p []byte) (n int, err error) {
	log.Printf("[CACHE] dirWriter.Write: BEFORE len=%d", len(p))
	n, err = d.f.Write(p)
	log.Printf("[CACHE] dirWriter.Write: AFTER n=%d err=%v", n, err)
	return n, err
}
```

---

## dirWriter.Complete (line 246)

```go
func (d *dirWriter) Complete() {
	log.Printf("[CACHE] dirWriter.Complete: setting complete=true")
	d.complete = true
}
```

---

## dirWriter.Close (line 250)

```go
func (d *dirWriter) Close() error {
	name := d.f.Name()
	log.Printf("[CACHE] dirWriter.Close: name=%s complete=%v dst=%s", name, d.complete, d.dst)
	log.Printf("[CACHE] dirWriter.Close: BEFORE d.f.Close")
	if err := d.f.Close(); err != nil {
		log.Printf("[CACHE] dirWriter.Close: AFTER d.f.Close err=%v", err)
		return fmt.Errorf("closing: %w", err)
	}
	log.Printf("[CACHE] dirWriter.Close: AFTER d.f.Close err=nil")
	if !d.complete {
		log.Printf("[CACHE] dirWriter.Close: BEFORE os.Remove (incomplete) name=%s", name)
		err := os.Remove(name)
		log.Printf("[CACHE] dirWriter.Close: AFTER os.Remove err=%v", err)
		return nil
	}
	log.Printf("[CACHE] dirWriter.Close: BEFORE os.Rename name=%s dst=%s", name, d.dst)
	if err := os.Rename(name, d.dst); err != nil {
		log.Printf("[CACHE] dirWriter.Close: AFTER os.Rename err=%v", err)
		return fmt.Errorf("renaming: %w", err)
	}
	log.Printf("[CACHE] dirWriter.Close: AFTER os.Rename err=nil")
	return nil
}
```

---

## memCache.get (line 281)

```go
func (m *memCache) get(ctx context.Context, key string) (*cacheEntry, error) {
	log.Printf("[MEM] get: key=%s", key)
	m.Lock()
	defer m.Unlock()

	for _, e := range m.entries {
		if e.key == key {
			e.access = time.Now()
			log.Printf("[MEM] get: found key=%s", key)
			return e, nil
		}
	}
	log.Printf("[MEM] get: not found key=%s", key)
	return nil, io.EOF
}
```

---

## memCache.Get (line 294)

```go
func (m *memCache) Get(ctx context.Context, key string) (*soci.TOC, error) {
	log.Printf("[MEM] Get: key=%s", key)
	e, err := m.get(ctx, key)
	if err != nil {
		log.Printf("[MEM] Get: err=%v", err)
		return nil, err
	}
	log.Printf("[MEM] Get: found")
	return e.toc, nil
}
```

---

## memCache.Put (line 303)

```go
func (m *memCache) Put(ctx context.Context, key string, toc *soci.TOC) error {
	log.Printf("[MEM] Put: key=%s size=%d", key, toc.Size)
	m.Lock()
	defer m.Unlock()
	if toc.Size > m.maxSize {
		log.Printf("[MEM] Put: too large, skipping")
		return nil
	}

	e := &cacheEntry{
		key:    key,
		toc:    toc,
		size:   toc.Size,
		access: time.Now(),
	}

	if len(m.entries) >= m.entryCap {
		min, idx := e.access, -1
		for i, e := range m.entries {
			if e.access.Before(min) {
				min = e.access
				idx = i
			}
		}
		log.Printf("[MEM] Put: evicting entry at idx=%d", idx)
		m.entries[idx] = e
		return nil
	}

	log.Printf("[MEM] Put: appending new entry")
	m.entries = append(m.entries, e)
	return nil
}
```

---

## memCache.New (line 335)

```go
func (m *memCache) New(ctx context.Context, key string) *cacheEntry {
	log.Printf("[MEM] New: key=%s", key)
	e := &cacheEntry{
		key:    key,
		access: time.Now(),
	}
	if len(m.entries) >= m.entryCap {
		min, idx := e.access, -1
		for i, e := range m.entries {
			if e.access.Before(min) {
				min = e.access
				idx = i
			}
		}
		log.Printf("[MEM] New: evicting entry at idx=%d", idx)
		m.entries[idx] = e
	} else {
		log.Printf("[MEM] New: appending")
		m.entries = append(m.entries, e)
	}
	return e
}
```

---

## memWriter.Write (line 360)

```go
func (w *memWriter) Write(p []byte) (n int, err error) {
	log.Printf("[MEM] memWriter.Write: len=%d", len(p))
	return w.buf.Write(p)
}
```

---

## memWriter.Close (line 364)

```go
func (w *memWriter) Close() (err error) {
	log.Printf("[MEM] memWriter.Close: bufLen=%d", w.buf.Len())
	w.entry.buffer = w.buf.Bytes()
	return nil
}
```

---

## memCache.Writer (line 369)

```go
func (m *memCache) Writer(ctx context.Context, key string) (io.WriteCloser, error) {
	log.Printf("[MEM] Writer: key=%s", key)
	e := m.New(ctx, key)
	mw := &memWriter{entry: e, buf: bytes.NewBuffer([]byte{})}
	return mw, nil
}
```

---

## memCache.Reader (line 375)

```go
func (m *memCache) Reader(ctx context.Context, key string) (io.ReadCloser, error) {
	log.Printf("[MEM] Reader: key=%s", key)
	e, err := m.get(ctx, key)
	if err != nil {
		log.Printf("[MEM] Reader: err=%v", err)
		return nil, err
	}
	log.Printf("[MEM] Reader: found, bufLen=%d", len(e.buffer))
	return io.NopCloser(bytes.NewReader(e.buffer)), nil
}
```

---

## memCache.RangeReader (line 383)

```go
func (m *memCache) RangeReader(ctx context.Context, key string, offset, length int64) (io.ReadCloser, error) {
	log.Printf("[MEM] RangeReader: key=%s offset=%d length=%d", key, offset, length)
	e, err := m.get(ctx, key)
	if err != nil {
		return nil, err
	}

	if offset == 0 && length == -1 {
		return m.Reader(ctx, key)
	}
	if e.buffer == nil || int64(len(e.buffer)) < offset+length+1 {
		log.Printf("[MEM] RangeReader: buffer too small")
		return nil, io.EOF
	}
	return io.NopCloser(bytes.NewReader(e.buffer[offset : offset+length])), nil
}
```

---

## memCache.Size (line 398)

```go
func (m *memCache) Size(ctx context.Context, key string) (int64, error) {
	log.Printf("[MEM] Size: key=%s", key)
	e, err := m.get(ctx, key)
	if err != nil {
		return -1, err
	}
	log.Printf("[MEM] Size: size=%d", len(e.buffer))
	return int64(len(e.buffer)), nil
}
```

---

## memCache.Delete (line 406)

```go
func (m *memCache) Delete(ctx context.Context, key string) error {
	log.Printf("[MEM] Delete: key=%s", key)
	m.Lock()
	defer m.Unlock()
	for i, e := range m.entries {
		if e.key == key {
			log.Printf("[MEM] Delete: found at idx=%d", i)
			m.entries = append(m.entries[:i], m.entries[i+1:]...)
			return nil
		}
	}
	log.Printf("[MEM] Delete: not found")
	return nil
}
```

---

## multiCache.Get (line 422)

```go
func (m *multiCache) Get(ctx context.Context, key string) (*soci.TOC, error) {
	log.Printf("[MULTI] Get: key=%s", key)
	for i, c := range m.caches {
		log.Printf("[MULTI] Get: trying cache[%d] %T", i, c)
		toc, err := c.Get(ctx, key)
		if err == nil {
			log.Printf("[MULTI] Get: hit in cache[%d]", i)
			for j := i - 1; j >= 0; j-- {
				cache := m.caches[j]
				log.Printf("[MULTI] Get: backfilling cache[%d] %T", j, cache)
				if err := cache.Put(ctx, key, toc); err != nil {
					log.Printf("[MULTI] Get: backfill err=%v", err)
				}
			}
			return toc, err
		} else {
			log.Printf("[MULTI] Get: miss in cache[%d] err=%v", i, err)
		}
	}
	log.Printf("[MULTI] Get: miss in all caches")
	return nil, io.EOF
}
```

---

## multiCache.Put (line 445)

```go
func (m *multiCache) Put(ctx context.Context, key string, toc *soci.TOC) error {
	log.Printf("[MULTI] Put: key=%s", key)
	errs := []error{}
	for i, c := range m.caches {
		log.Printf("[MULTI] Put: putting in cache[%d] %T", i, c)
		err := c.Put(ctx, key, toc)
		if err != nil {
			log.Printf("[MULTI] Put: cache[%d] err=%v", i, err)
			errs = append(errs, err)
		}
	}
	return Join(errs...)
}
```

---

## multiCache.Writer (line 457)

```go
func (m *multiCache) Writer(ctx context.Context, key string) (io.WriteCloser, error) {
	log.Printf("[MULTI] Writer: key=%s", key)
	writers := []io.WriteCloser{}
	for i, c := range m.caches {
		log.Printf("[MULTI] Writer: getting writer from cache[%d] %T", i, c)
		w, err := c.Writer(ctx, key)
		if err != nil {
			log.Printf("[MULTI] Writer: cache[%d] err=%v", i, err)
			return nil, err
		}
		writers = append(writers, w)
	}
	log.Printf("[MULTI] Writer: returning MultiWriter with %d writers", len(writers))
	return MultiWriter(writers...), nil
}
```

---

## multiCache.Reader (line 470)

```go
func (m *multiCache) Reader(ctx context.Context, key string) (io.ReadCloser, error) {
	log.Printf("[MULTI] Reader: key=%s", key)
	for i, c := range m.caches {
		log.Printf("[MULTI] Reader: trying cache[%d] %T", i, c)
		rc, err := c.Reader(ctx, key)
		if err == nil {
			log.Printf("[MULTI] Reader: hit in cache[%d]", i)
			return rc, nil
		} else {
			log.Printf("[MULTI] Reader: miss in cache[%d] err=%v", i, err)
		}
	}
	log.Printf("[MULTI] Reader: miss in all caches")
	return nil, io.EOF
}
```

---

## multiCache.RangeReader (line 484)

```go
func (m *multiCache) RangeReader(ctx context.Context, key string, offset, length int64) (io.ReadCloser, error) {
	log.Printf("[MULTI] RangeReader: key=%s offset=%d length=%d", key, offset, length)
	for i, c := range m.caches {
		var (
			rc  io.ReadCloser
			err error
		)
		if offset == 0 && length == -1 {
			rc, err = c.Reader(ctx, key)
		} else {
			rc, err = c.RangeReader(ctx, key, offset, length)
		}
		if err == nil {
			log.Printf("[MULTI] RangeReader: hit in cache[%d]", i)
			return rc, nil
		} else {
			log.Printf("[MULTI] RangeReader: miss in cache[%d] err=%v", i, err)
		}
	}
	log.Printf("[MULTI] RangeReader: miss in all caches")
	return nil, io.EOF
}
```

---

## multiCache.Size (line 506)

```go
func (m *multiCache) Size(ctx context.Context, key string) (int64, error) {
	log.Printf("[MULTI] Size: key=%s", key)
	for i, c := range m.caches {
		log.Printf("[MULTI] Size: trying cache[%d] %T", i, c)
		sz, err := c.Size(ctx, key)
		if err == nil {
			log.Printf("[MULTI] Size: hit in cache[%d] size=%d", i, sz)
			return sz, nil
		} else {
			log.Printf("[MULTI] Size: miss in cache[%d] err=%v", i, err)
		}
	}
	log.Printf("[MULTI] Size: miss in all caches")
	return -1, io.EOF
}
```

---

## multiCache.Delete (line 519)

```go
func (m *multiCache) Delete(ctx context.Context, key string) error {
	log.Printf("[MULTI] Delete: key=%s", key)
	errs := []error{}
	for i, c := range m.caches {
		log.Printf("[MULTI] Delete: deleting from cache[%d] %T", i, c)
		if err := c.Delete(ctx, key); err != nil {
			log.Printf("[MULTI] Delete: cache[%d] err=%v", i, err)
			errs = append(errs, err)
		}
	}
	return Join(errs...)
}
```

---

## multiWriter.Write (line 533)

```go
func (t *multiWriter) Write(p []byte) (n int, err error) {
	log.Printf("[MULTI] multiWriter.Write: len=%d", len(p))
	for i, w := range t.writers {
		log.Printf("[MULTI] multiWriter.Write: writing to writer[%d]", i)
		n, err = w.Write(p)
		if err != nil {
			log.Printf("[MULTI] multiWriter.Write: writer[%d] err=%v", i, err)
			return
		}
		if n != len(p) {
			err = io.ErrShortWrite
			log.Printf("[MULTI] multiWriter.Write: writer[%d] short write", i)
			return
		}
	}
	return len(p), nil
}
```

---

## multiWriter.Complete (line 547)

```go
func (t *multiWriter) Complete() {
	log.Printf("[MULTI] multiWriter.Complete: numWriters=%d", len(t.writers))
	for i, w := range t.writers {
		if cw, ok := w.(interface{ Complete() }); ok {
			log.Printf("[MULTI] multiWriter.Complete: calling Complete on writer[%d]", i)
			cw.Complete()
		}
	}
}
```

---

## multiWriter.Close (line 555)

```go
func (t *multiWriter) Close() error {
	log.Printf("[MULTI] multiWriter.Close: numWriters=%d", len(t.writers))
	errs := []error{}
	for i, w := range t.writers {
		log.Printf("[MULTI] multiWriter.Close: closing writer[%d]", i)
		if err := w.Close(); err != nil {
			log.Printf("[MULTI] multiWriter.Close: writer[%d] err=%v", i, err)
			errs = append(errs, err)
		}
	}
	return Join(errs...)
}
```

---

## buildGcsCache (line 619)

```go
func buildGcsCache(bucket string) (cache, error) {
	log.Printf("[BUILD] buildGcsCache: bucket=%s", bucket)
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Printf("[BUILD] buildGcsCache: NewClient err=%v", err)
		return nil, err
	}
	bkt := client.Bucket(strings.TrimPrefix(bucket, "gs://"))
	log.Printf("[BUILD] buildGcsCache: created bucket handle")
	return &gcsCache{client, bkt}, nil
}
```

---

## buildTocCache (line 629)

```go
func buildTocCache() cache {
	log.Printf("[BUILD] buildTocCache: START")
	mc := &memCache{
		maxSize:  50 * (1 << 20),
		entryCap: 50,
	}
	log.Printf("[BUILD] buildTocCache: created memCache maxSize=%d entryCap=%d", mc.maxSize, mc.entryCap)

	caches := []cache{mc}
	caches = append(caches, &dirCache{dir: "/cache"})
	log.Printf("[BUILD] buildTocCache: added dirCache dir=/cache")

	log.Printf("[BUILD] buildTocCache: returning multiCache with %d caches", len(caches))
	return &multiCache{caches}
}
```

---

## buildIndexCache (line 643)

```go
func buildIndexCache() cache {
	log.Printf("[BUILD] buildIndexCache: START")
	caches := []cache{}
	caches = append(caches, &dirCache{dir: "/cache"})
	log.Printf("[BUILD] buildIndexCache: added dirCache dir=/cache")
	log.Printf("[BUILD] buildIndexCache: returning multiCache with %d caches", len(caches))
	return &multiCache{caches}
}
```
