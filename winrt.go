package winrt

// event
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Foundation.TypedEventHandler`2
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Foundation.EventRegistrationToken

// vector
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Foundation.Collections.IVector`1
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Foundation.Collections.IVectorView`1

// TimeSpan
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Foundation.TimeSpan
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Foundation.DateTime

// smtc
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Media.SoundLevel
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Media.MediaPlaybackStatus
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Media.MediaPlaybackAutoRepeatMode
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Media.SystemMediaTransportControls
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Media.SystemMediaTransportControlsTimelineProperties
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Media.MediaPlaybackType
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Media.MusicDisplayProperties
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Media.VideoDisplayProperties
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Media.ImageDisplayProperties
//go:generate go run github.com/waylyrics/winrt-go/cmd/winrt-go-gen -debug -class Windows.Media.SystemMediaTransportControlsDisplayUpdater -method-filter !CopyFromFileAsync -method-filter !get_Thumbnail -method-filter !put_Thumbnail
