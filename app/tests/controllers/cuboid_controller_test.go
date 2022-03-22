package controllers_test

import (
	. "cuboid-challenge/app/models"
	"cuboid-challenge/app/tests/testutils"
	"fmt"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Cuboid Controller", func() {
	testutils.LoadEnv()
	testutils.ConnectDB()
	testutils.ClearDB()

	AfterEach(func() {
		testutils.ClearDB()
	})

	var w *httptest.ResponseRecorder
	var bag *Bag

	BeforeEach(func() {
		bag = &Bag{
			Title:  "A bag",
			Volume: 5,
			Cuboids: []Cuboid{
				{Width: 1, Height: 1, Depth: 1},
				{Width: 1, Height: 1, Depth: 2},
			},
		}
		testutils.AddRecords(bag)
	})

	Describe("List", func() {
		BeforeEach(func() {
			w = testutils.MockRequest(http.MethodGet, "/cuboids", nil)
		})

		It("Response HTTP status code 200", func() {
			Expect(w.Code).To(Equal(200))
		})

		It("Fetch all cuboids", func() {
			l, _ := testutils.DeserializeList(w.Body.String())
			Expect(len(l)).To(Equal(2))
			for _, m := range l {
				Expect(m["width"]).ToNot(BeNil())
				Expect(m["height"]).ToNot(BeNil())
				Expect(m["depth"]).ToNot(BeNil())
				Expect(m["volume"]).ToNot(BeNil())
				Expect(m["bagId"]).To(BeEquivalentTo(bag.ID))
			}
		})
	})

	Describe("Get", func() {
		var cuboidID uint

		JustBeforeEach(func() {
			w = testutils.MockRequest(http.MethodGet, "/cuboids/"+fmt.Sprintf("%v", cuboidID), nil)
		})

		Context("When the cuboid is present", func() {
			BeforeEach(func() {
				cuboidID = bag.Cuboids[0].ID
			})

			It("Response HTTP status code 200", func() {
				Expect(w.Code).To(Equal(200))
			})

			It("Get the cuboid", func() {
				m, _ := testutils.Deserialize(w.Body.String())
				Expect(m["id"]).To(BeEquivalentTo(bag.Cuboids[0].ID))
				Expect(m["width"]).ToNot(BeNil())
				Expect(m["height"]).ToNot(BeNil())
				Expect(m["depth"]).ToNot(BeNil())
				Expect(m["volume"]).ToNot(BeNil())
				Expect(m["bagId"]).ToNot(BeNil())
			})
		})

		Context("When cuboid is not present", func() {
			BeforeEach(func() {
				cuboidID = 9999
			})

			It("Response HTTP status code 404", func() {
				Expect(w.Code).To(Equal(404))
			})
		})
	})

	Describe("Create", func() {
		cuboidPayload := map[string]interface{}{}

		BeforeEach(func() {
			cuboidPayload = map[string]interface{}{
				"width":  1,
				"height": 1,
				"depth":  1,
				"bagId":  bag.ID,
			}
		})

		JustBeforeEach(func() {
			body, _ := testutils.SerializeToString(cuboidPayload)
			w = testutils.MockRequest(http.MethodPost, "/cuboids", &body)
		})

		It("Response HTTP status code 201", func() {
			Expect(w.Code).To(Equal(201))
		})

		It("Returns the created cuboid", func() {
			m, _ := testutils.Deserialize(w.Body.String())
			Expect(m["width"]).ToNot(BeNil())
			Expect(m["height"]).ToNot(BeNil())
			Expect(m["depth"]).ToNot(BeNil())
			Expect(m["volume"]).ToNot(BeNil())
			Expect(m["bagId"]).To(BeEquivalentTo(bag.ID))
		})

		Context("When cuboid does not fit into the bag", func() {
			BeforeEach(func() {
				cuboidPayload["width"] = 3
			})

			It("Does not create the Cuboid", func() {
				Expect(w.Code).To(Equal(400))
				m, _ := testutils.Deserialize(w.Body.String())
				Expect(m["error"]).To(Equal("Insufficient capacity in bag"))
			})
		})

		Context("When the bag is disabled", func() {
			BeforeEach(func() {
				bag.SetDisabled(true)
				bag.Cuboids = []Cuboid{}
				testutils.UpdateRecords(&bag)
			})

			It("Does not admit new cuboids", func() {
				Expect(w.Code).To(Equal(400))
				m, _ := testutils.Deserialize(w.Body.String())
				Expect(m["error"]).To(Equal("Bag is disabled"))
			})
		})
	})

	// DO NOT modify the tests ABOVE
	// IMPLEMENT the tests BELOW

	Describe("Update", func() {
		var newCuboid *Cuboid

		BeforeEach(func() {
			newCuboid = &Cuboid{
				Width:  2,
				Height: 2,
				Depth:  2,
				BagID:  bag.ID,
				Bag:    bag,
			}
			testutils.AddRecords(newCuboid)
		})

		JustBeforeEach(func() {
			body, _ := testutils.SerializeToString(newCuboid)
			w = testutils.MockRequest(http.MethodPut, "/cuboids/"+fmt.Sprintf("%v", newCuboid.ID), &body)
		})

		It("Response HTTP status code 200", func() {
			Expect(w.Code).To(Equal(200))
		})

		It("Returns the updated cuboid", func() {
			m, _ := testutils.Deserialize(w.Body.String())
			Expect(m["width"]).To(BeEquivalentTo(newCuboid.Width))
			Expect(m["height"]).To(BeEquivalentTo(newCuboid.Height))
			Expect(m["depth"]).To(BeEquivalentTo(newCuboid.Depth))
			Expect(m["bagId"]).To(BeEquivalentTo(newCuboid.BagID))
		})

		Context("When cuboid does not fit into the bag", func() {
			BeforeEach(func() {
				newCuboid.Width = 999
				newCuboid.Height = 999
				newCuboid.Depth = 999
			})

			It("Response HTTP status code 400", func() {
				Expect(w.Code).To(Equal(400))

			})

			It("Response a JSON with error message 'Insufficient capacity in bag'", func() {
				m, _ := testutils.Deserialize(w.Body.String())
				Expect(m["error"]).To(Equal("Insufficient capacity in bag"))
			})
		})

		Context("When cuboid is not present", func() {
			BeforeEach(func() {
				newCuboid.ID = 9999
			})

			It("Response HTTP status code 404", func() {
				Expect(w.Code).To(Equal(404))
			})
		})
	})

	Describe("Delete", func() {
		var cuboidID uint
		var cuboid Cuboid

		BeforeEach(func() {
			cuboidID = bag.Cuboids[0].ID
			cuboid = bag.Cuboids[0]
		})

		JustBeforeEach(func() {
			w = testutils.MockRequest(http.MethodDelete, "/cuboids/"+fmt.Sprintf("%v", cuboidID), nil)
		})

		Context("When the cuboid is present", func() {
			It("Response HTTP status code 200", func() {
				Expect(w.Code).To(Equal(200))
			})

			It("Remove the cuboid", func() {
				m, _ := testutils.Deserialize(w.Body.String())
				Expect(m["width"]).To(BeEquivalentTo(cuboid.Width))
				Expect(m["height"]).To(BeEquivalentTo(cuboid.Height))
				Expect(m["depth"]).To(BeEquivalentTo(cuboid.Depth))
				Expect(m["bagId"]).To(BeEquivalentTo(cuboid.BagID))
			})
		})

		Context("When cuboid is not present", func() {
			BeforeEach(func() {
				cuboidID = 9999
			})

			It("Response HTTP status code 400", func() {
				Expect(w.Code).To(Equal(404))
			})
		})
	})
})
