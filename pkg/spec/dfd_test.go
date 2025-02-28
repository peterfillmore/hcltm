package spec

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
)

func dfdTm() *Threatmodel {
	tm := &Threatmodel{
		Name:   "test",
		Author: "x",
		DataFlowDiagram: &DataFlowDiagram{
			Processes: []*DfdProcess{
				&DfdProcess{
					Name: "proc1",
				},
			},
		},
	}

	return tm
}

func fullDfdTm() *Threatmodel {

	tm := &Threatmodel{
		Name:   "test",
		Author: "x",
		DataFlowDiagram: &DataFlowDiagram{
			Processes: []*DfdProcess{
				&DfdProcess{
					Name: "proc1",
				},
				&DfdProcess{
					Name:      "proc2",
					TrustZone: "zone1",
				},
			},
			DataStores: []*DfdData{
				&DfdData{
					Name: "data1",
				},
				&DfdData{
					Name:      "data2",
					TrustZone: "zone2",
				},
			},
			ExternalElements: []*DfdExternal{
				&DfdExternal{
					Name: "external1",
				},
				&DfdExternal{
					Name:      "external2",
					TrustZone: "zone3",
				},
			},
			Flows: []*DfdFlow{
				&DfdFlow{
					Name: "flow",
					From: "proc1",
					To:   "data1",
				},
				&DfdFlow{
					Name: "flow",
					From: "external1",
					To:   "proc1",
				},
				&DfdFlow{
					Name: "flow",
					From: "data1",
					To:   "external1",
				},
			},
		},
	}

	return tm

}

func fullDfdTm2() *Threatmodel {

	tm := &Threatmodel{
		Name:   "test",
		Author: "x",
		DataFlowDiagram: &DataFlowDiagram{
			TrustZones: []*DfdTrustZone{
				&DfdTrustZone{
					Name: "zone1",
					Processes: []*DfdProcess{
						&DfdProcess{
							Name:      "proc2",
							TrustZone: "zone1",
						},
						&DfdProcess{
							Name: "proc9",
						},
					},
					DataStores: []*DfdData{
						&DfdData{
							Name: "new_data",
						},
					},
					ExternalElements: []*DfdExternal{
						&DfdExternal{
							Name: "ee5",
						},
					},
				},
			},
			Processes: []*DfdProcess{
				&DfdProcess{
					Name: "proc1",
				},
			},
			DataStores: []*DfdData{
				&DfdData{
					Name: "data1",
				},
				&DfdData{
					Name:      "data2",
					TrustZone: "zone2",
				},
			},
			ExternalElements: []*DfdExternal{
				&DfdExternal{
					Name: "external1",
				},
				&DfdExternal{
					Name:      "external2",
					TrustZone: "zone3",
				},
			},
			Flows: []*DfdFlow{
				&DfdFlow{
					Name: "flow",
					From: "proc1",
					To:   "data1",
				},
				&DfdFlow{
					Name: "flow",
					From: "external1",
					To:   "proc1",
				},
				&DfdFlow{
					Name: "flow",
					From: "data1",
					To:   "external1",
				},
			},
		},
	}

	return tm

}

func TestDfdPngGenerate(t *testing.T) {
	// tm := dfdTm()
	//
	// fulltm := fullDfdTm()

	cases := []struct {
		name        string
		tm          *Threatmodel
		exp         string
		errorthrown bool
	}{
		{
			"valid_dfd",
			dfdTm(),
			"",
			false,
		},
		{
			"valid_full_dfd",
			fullDfdTm(),
			"",
			false,
		},
		{
			"valid_full_dfd2",
			fullDfdTm2(),
			"",
			false,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()

			d, err := ioutil.TempDir("", "")
			if err != nil {
				t.Fatalf("Error creating tmp dir: %s", err)
			}
			defer os.RemoveAll(d)

			err = tc.tm.GenerateDfdPng(fmt.Sprintf("%s/out.png", d))

			if err != nil {
				if !strings.Contains(err.Error(), tc.exp) {
					t.Errorf("%s: Error rendering png: %s", tc.name, err)
				}
			} else {
				if tc.errorthrown {
					t.Errorf("%s: an error was thrown when it shoulnd't have", tc.name)
				} else {

					// at this point we should have a legitimate png to
					// test

					f, err := os.Open(fmt.Sprintf("%s/out.png", d))
					if err != nil {
						t.Fatalf("%s: Error opening png: %s", tc.name, err)
					}

					buffer := make([]byte, 512)
					_, err = f.Read(buffer)
					if err != nil {
						t.Fatalf("%s: Error reading png: %s", tc.name, err)
					}

					if http.DetectContentType(buffer) != "image/png" {
						t.Errorf("%s: The output file isn't a png, it's '%s'", tc.name, http.DetectContentType(buffer))
					}
				}
			}

		})
	}
}

func TestDfdSvgGenerate(t *testing.T) {
	// tm := dfdTm()
	//
	// fulltm := fullDfdTm()

	cases := []struct {
		name        string
		tm          *Threatmodel
		exp         string
		errorthrown bool
	}{
		{
			"valid_dfd",
			dfdTm(),
			"",
			false,
		},
		{
			"valid_full_dfd",
			fullDfdTm(),
			"",
			false,
		},
		{
			"valid_full_dfd2",
			fullDfdTm2(),
			"",
			false,
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()

			d, err := ioutil.TempDir("", "")
			if err != nil {
				t.Fatalf("Error creating tmp dir: %s", err)
			}
			defer os.RemoveAll(d)

			err = tc.tm.GenerateDfdSvg(fmt.Sprintf("%s/out.svg", d))

			if err != nil {
				if !strings.Contains(err.Error(), tc.exp) {
					t.Errorf("%s: Error rendering svg: %s", tc.name, err)
				}
			} else {
				if tc.errorthrown {
					t.Errorf("%s: an error was thrown when it shouldn't have", tc.name)
				} else {

					// at this point we should have a legitimate svg to
					// test

					f, err := os.Open(fmt.Sprintf("%s/out.svg", d))
					if err != nil {
						t.Fatalf("%s: Error opening svg: %s", tc.name, err)
					}

					buffer := make([]byte, 512)
					_, err = f.Read(buffer)
					if err != nil {
						t.Fatalf("%s: Error reading svg: %s", tc.name, err)
					}

					if http.DetectContentType(buffer) != "text/xml; charset=utf-8" {
						t.Errorf("%s: The output file isn't a svg, it's '%s'", tc.name, http.DetectContentType(buffer))
					}
				}
			}

		})
	}
}